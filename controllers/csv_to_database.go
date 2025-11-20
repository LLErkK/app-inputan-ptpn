package controllers

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"app-inputan-ptpn/config"
	"app-inputan-ptpn/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// parseNumber: robust untuk "1.234,56", "1,234.56", "(123)", "12,5%" dll.
func parseNumber(raw string) (float64, error) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return 0, fmt.Errorf("empty")
	}

	// normalize spaces incl NBSP
	s = strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return ' '
		}
		return r
	}, s)
	s = strings.ReplaceAll(s, "\u00A0", "")
	s = strings.TrimSpace(s)

	negative := false
	if strings.HasPrefix(s, "(") && strings.HasSuffix(s, ")") {
		negative = true
		s = strings.TrimPrefix(strings.TrimSuffix(s, ")"), "(")
		s = strings.TrimSpace(s)
	}

	hasPercent := false
	if strings.HasSuffix(s, "%") {
		hasPercent = true
		s = strings.TrimSuffix(s, "%")
		s = strings.TrimSpace(s)
	}

	re := regexp.MustCompile(`[^0-9\.,\-]`)
	s = re.ReplaceAllString(s, "")

	if s == "" || s == "-" {
		return 0, fmt.Errorf("no-number")
	}

	if strings.Contains(s, ".") && strings.Contains(s, ",") {
		if strings.LastIndex(s, ",") > strings.LastIndex(s, ".") {
			s = strings.ReplaceAll(s, ".", "")
			s = strings.ReplaceAll(s, ",", ".")
		} else {
			s = strings.ReplaceAll(s, ",", "")
		}
	} else if strings.Contains(s, ",") && !strings.Contains(s, ".") {
		s = strings.ReplaceAll(s, ".", "")
		s = strings.ReplaceAll(s, ",", ".")
	} else {
		if strings.Count(s, ".") > 1 {
			lastDot := strings.LastIndex(s, ".")
			after := len(s) - lastDot - 1
			if after == 3 {
				s = strings.ReplaceAll(s, ".", "")
			}
		}
	}

	if negative {
		s = "-" + s
	}

	num, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		return 0, fmt.Errorf("parse error: %v (%s)", err, s)
	}

	_ = hasPercent
	return num, nil
}

// findHeaderRowAndBaseIndex: cari row & kolom yang mengandung "TAHUN TANAM"
func findHeaderRowAndBaseIndex(rows [][]string, maxScan int) (int, int) {
	limit := maxScan
	if len(rows) < limit {
		limit = len(rows)
	}
	for i := 0; i < limit; i++ {
		row := rows[i]
		for j, cell := range row {
			clean := strings.ToUpper(strings.TrimSpace(cell))
			if clean == "" {
				continue
			}
			c := strings.ReplaceAll(clean, "_", " ")
			c = strings.ReplaceAll(c, "\"", "")
			c = strings.TrimSpace(c)

			if strings.Contains(c, "TAHUN TANAM") || c == "TAHUN" || c == "TAHUN_TANAM" {
				return i, j
			}
			if c == "NIK" {
				return i, j - 1
			}
		}
	}
	return -1, -1
}

// detectTipeProduksi: cari di SEMUA kolom
func detectTipeProduksi(row []string, baseIdx int) string {
	for i := 0; i < len(row); i++ {
		cell := strings.ToUpper(strings.TrimSpace(row[i]))
		if cell == "" {
			continue
		}

		// Urutan pengecekan dari yang paling spesifik ke general
		if strings.Contains(cell, "REKAPITULASI") {
			return "REKAPITULASI"
		}
		if strings.Contains(cell, "PRODUKSI TETES LANJUT") || strings.Contains(cell, "TETES LANJUT") {
			return "PRODUKSI TETES LANJUT"
		}
		if strings.Contains(cell, "PRODUKSI BORONG EKSTERNAL") || strings.Contains(cell, "BORONG EKSTERNAL") {
			return "PRODUKSI BORONG EKSTERNAL"
		}
		if strings.Contains(cell, "PRODUKSI BORONG INTERNAL") || strings.Contains(cell, "BORONG INTERNAL") {
			return "PRODUKSI BORONG INTERNAL"
		}
		if strings.Contains(cell, "PRODUKSI BORONG MINGGU") || strings.Contains(cell, "BORONG MINGGU") {
			return "PRODUKSI BORONG MINGGU"
		}
		if strings.Contains(cell, "PRODUKSI BAKU BORONG") || strings.Contains(cell, "BAKU BORONG") {
			return "PRODUKSI BAKU BORONG"
		}
		if strings.Contains(cell, "PRODUKSI BAKU") {
			return "PRODUKSI BAKU"
		}
	}
	return ""
}

// isTipeProduksiRow: lebih robust dengan cek NIK dan TAHUN TANAM kosong
func isTipeProduksiRow(row []string, baseIdx int) bool {
	tipe := detectTipeProduksi(row, baseIdx)
	if tipe == "" {
		return false
	}

	tahunTanam := ""
	if baseIdx >= 0 && baseIdx < len(row) {
		tahunTanam = strings.TrimSpace(row[baseIdx])
	}

	nik := ""
	if baseIdx+1 >= 0 && baseIdx+1 < len(row) {
		nik = strings.TrimSpace(row[baseIdx+1])
	}

	if tahunTanam == "" && nik == "" {
		return true
	}

	if tahunTanam != "" {
		if regexp.MustCompile(`^\d{4}$`).MatchString(tahunTanam) {
			return false
		}
	}

	return true
}

// isLikelySummaryRow: tambah pengecekan untuk baris kategori produksi
func isLikelySummaryRow(row []string, baseIdx int) bool {
	if isTipeProduksiRow(row, baseIdx) {
		return false
	}

	checkIdx := baseIdx + 1
	if checkIdx < 0 || checkIdx >= len(row) {
		return false
	}
	cell := strings.ToUpper(strings.TrimSpace(row[checkIdx]))

	if cell == "" {
		c0 := strings.ToUpper(strings.TrimSpace(row[baseIdx]))
		if strings.Contains(c0, "JUMLAH") || strings.Contains(c0, "SELISIH") ||
			strings.Contains(c0, "%") || strings.Contains(c0, "K3") ||
			strings.HasPrefix(c0, "√") || c0 == "OW" {
			return true
		}
		return false
	}

	summaryKeys := []string{"JUMLAH", "SELISIH", "TOTAL", "K3", "%",
		"JUMLAH PABRIK", "JUMLAH KEBUN", "RATA",
		"RATA-RATA", "REKAPITULASI", "OW", "PROD. OW"}
	for _, k := range summaryKeys {
		if strings.Contains(cell, k) {
			return true
		}
	}

	if regexp.MustCompile(`^[A-Z ]+$`).MatchString(cell) && len(cell) <= 6 &&
		!regexp.MustCompile(`[0-9]`).MatchString(cell) {
		return true
	}
	return false
}

// isValidDataRow: validasi minimal baris data
func isValidDataRow(row []string, baseIdx int) bool {
	tyIdx := baseIdx
	if tyIdx < 0 || tyIdx >= len(row) {
		return false
	}
	ty := strings.TrimSpace(row[tyIdx])
	if ty == "" {
		return false
	}
	tyDigits := regexp.MustCompile(`\d{4}`)
	if !tyDigits.MatchString(ty) {
		return false
	}
	yearStr := tyDigits.FindString(ty)
	if y, err := strconv.Atoi(yearStr); err != nil || y < 1900 || y > 2100 {
		return false
	}

	nikIdx := baseIdx + 1
	if nikIdx < 0 || nikIdx >= len(row) {
		return false
	}
	nik := strings.TrimSpace(row[nikIdx])
	if nik == "" {
		return false
	}
	nikClean := strings.ReplaceAll(nik, ".", "")
	nikClean = strings.ReplaceAll(nikClean, " ", "")
	if !regexp.MustCompile(`\d`).MatchString(nikClean) {
		return false
	}
	digits := regexp.MustCompile(`\d+`).FindString(nikClean)
	if len(digits) < 4 {
		return false
	}

	return true
}

// hasValidHKO: cek apakah HKO ada nilainya
func hasValidHKO(row []string, baseIdx int) bool {
	getInt := func(idx int) int {
		if idx < 0 || idx >= len(row) {
			return 0
		}
		v := strings.TrimSpace(strings.ReplaceAll(row[idx], "\"", ""))
		v = strings.ReplaceAll(v, "-", "")
		v = strings.ReplaceAll(v, "—", "")
		v = strings.TrimSpace(v)

		if v == "" {
			return 0
		}
		if i, err := strconv.Atoi(strings.ReplaceAll(v, ".", "")); err == nil {
			return i
		}
		if f, err := parseNumber(v); err == nil {
			return int(f)
		}
		return 0
	}

	hkoHariIni := getInt(baseIdx + 3)
	hkoSampaiHariIni := getInt(baseIdx + 4)

	if hkoHariIni == 0 && hkoSampaiHariIni == 0 {
		return false
	}

	return true
}

// mapRowRelative: mapping urut sesuai struktur model Rekap
func mapRowRelative(row []string, baseIdx int, tanggal time.Time, tipeProduksi string, afdeling string, idMaster uint64) (*models.Rekap, error) {
	rekap := &models.Rekap{}

	getStr := func(idx int) string {
		if idx < 0 || idx >= len(row) {
			return ""
		}
		return strings.TrimSpace(strings.ReplaceAll(row[idx], "\"", ""))
	}

	getFloat := func(idx int) float64 {
		if idx < 0 || idx >= len(row) {
			return 0
		}
		v := strings.TrimSpace(strings.ReplaceAll(row[idx], "\"", ""))
		if v == "" || v == "-" || v == "—" {
			return 0
		}
		if n, err := parseNumber(v); err == nil {
			return n
		}
		return 0
	}

	getInt := func(idx int) int {
		if idx < 0 || idx >= len(row) {
			return 0
		}
		v := strings.TrimSpace(strings.ReplaceAll(row[idx], "\"", ""))
		if v == "" {
			return 0
		}
		if i, err := strconv.Atoi(strings.ReplaceAll(v, ".", "")); err == nil {
			return i
		}
		if f, err := parseNumber(v); err == nil {
			return int(f)
		}
		return 0
	}

	rekap.Tanggal = tanggal
	rekap.TipeProduksi = tipeProduksi
	rekap.TahunTanam = getStr(baseIdx + 0)
	rekap.NIK = getStr(baseIdx + 1)
	rekap.Mandor = getStr(baseIdx + 2)
	rekap.HKOHariIni = getInt(baseIdx + 3)
	rekap.HKOSampaiHariIni = getInt(baseIdx + 4)
	rekap.HariIniBasahLatekKebun = getFloat(baseIdx + 5)
	rekap.HariIniBasahLatekPabrik = getFloat(baseIdx + 6)
	rekap.HariIniBasahLatekPersen = getFloat(baseIdx + 7)
	rekap.HariIniBasahLumpKebun = getFloat(baseIdx + 8)
	rekap.HariIniBasahLumpPabrik = getFloat(baseIdx + 9)
	rekap.HariIniBasahLumpPersen = getFloat(baseIdx + 10)
	rekap.HariIniK3Sheet = getFloat(baseIdx + 11)
	rekap.HariIniKeringSheet = getFloat(baseIdx + 12)
	rekap.HariIniKeringBrCr = getFloat(baseIdx + 13)
	rekap.HariIniKeringJumlah = getFloat(baseIdx + 14)
	rekap.SampaiHariIniBasahLatekKebun = getFloat(baseIdx + 15)
	rekap.SampaiHariIniBasahLatekPabrik = getFloat(baseIdx + 16)
	rekap.SampaiHariIniBasahLatekPersen = getFloat(baseIdx + 17)
	rekap.SampaiHariIniBasahLumpKebun = getFloat(baseIdx + 18)
	rekap.SampaiHariIniBasahLumpPabrik = getFloat(baseIdx + 19)
	rekap.SampaiHariIniBasahLumpPersen = getFloat(baseIdx + 20)
	rekap.SampaiHariIniK3Sheet = getFloat(baseIdx + 21)
	rekap.SampaiHariIniKeringSheet = getFloat(baseIdx + 22)
	rekap.SampaiHariIniKeringBrCr = getFloat(baseIdx + 23)
	rekap.SampaiHariIniKeringJumlah = getFloat(baseIdx + 24)
	rekap.ProduksiPerTaperHariIni = getFloat(baseIdx + 25)
	rekap.ProduksiPerTaperSampaiHariIni = getFloat(baseIdx + 26)
	rekap.Afdeling = afdeling
	rekap.IdMaster = idMaster
	rekap.TotalProduksi = rekap.HariIniKeringSheet + rekap.HariIniBasahLumpPabrik

	return rekap, nil
}

// saveBatchRekap: save multiple rekap records in batch using upsert
func saveBatchRekap(db *gorm.DB, rekaps []*models.Rekap) error {
	if len(rekaps) == 0 {
		return nil
	}

	// Use batch insert with ON CONFLICT clause for upsert behavior
	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "tanggal"},
			{Name: "tipe_produksi"},
			{Name: "nik"},
			{Name: "mandor"},
			{Name: "tahun_tanam"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"hko_hari_ini", "hko_sampai_hari_ini",
			"hari_ini_basah_latek_kebun", "hari_ini_basah_latek_pabrik", "hari_ini_basah_latek_persen",
			"hari_ini_basah_lump_kebun", "hari_ini_basah_lump_pabrik", "hari_ini_basah_lump_persen",
			"hari_ini_k3_sheet", "hari_ini_kering_sheet", "hari_ini_kering_br_cr", "hari_ini_kering_jumlah",
			"sampai_hari_ini_basah_latek_kebun", "sampai_hari_ini_basah_latek_pabrik", "sampai_hari_ini_basah_latek_persen",
			"sampai_hari_ini_basah_lump_kebun", "sampai_hari_ini_basah_lump_pabrik", "sampai_hari_ini_basah_lump_persen",
			"sampai_hari_ini_k3_sheet", "sampai_hari_ini_kering_sheet", "sampai_hari_ini_kering_br_cr", "sampai_hari_ini_kering_jumlah",
			"produksi_per_taper_hari_ini", "produksi_per_taper_sampai_hari_ini",
			"afdeling", "id_master", "updated_at",
		}),
	}).CreateInBatches(rekaps, 100).Error
}

// processCSVFileAutoBaseWithFilter: optimized with batch processing
func processCSVFileAutoBaseWithFilter(db *gorm.DB, path string, tanggal time.Time, afdeling string, idMaster uint64) (int, int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, 0, fmt.Errorf("gagal buka file %s: %w", path, err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = -1
	r.LazyQuotes = true
	r.ReuseRecord = true // Reuse memory

	rows, err := r.ReadAll()
	if err != nil {
		return 0, 0, fmt.Errorf("gagal baca csv %s: %w", path, err)
	}
	if len(rows) == 0 {
		return 0, 0, fmt.Errorf("file kosong: %s", path)
	}

	headerRow, baseIdx := findHeaderRowAndBaseIndex(rows, 30)
	if headerRow == -1 || baseIdx == -1 {
		return 0, 0, fmt.Errorf("tidak menemukan header 'TAHUN TANAM' atau 'NIK' di %s", path)
	}

	fmt.Printf("DEBUG: Header found at row %d, baseIdx %d\n", headerRow, baseIdx)

	start := headerRow + 3
	saved, failed := 0, 0
	currentTipeProduksi := "PRODUKSI BAKU"

	// Batch processing
	const batchSize = 100
	rekaps := make([]*models.Rekap, 0, batchSize)
	var mu sync.Mutex

	for i := start; i < len(rows); i++ {
		row := rows[i]

		// skip empty rows
		nonEmpty := false
		for _, c := range row {
			if strings.TrimSpace(c) != "" {
				nonEmpty = true
				break
			}
		}
		if !nonEmpty {
			continue
		}

		if isTipeProduksiRow(row, baseIdx) {
			newTipe := detectTipeProduksi(row, baseIdx)
			if newTipe != "" {
				// Save current batch before changing type
				if len(rekaps) > 0 {
					mu.Lock()
					if err := saveBatchRekap(db, rekaps); err != nil {
						fmt.Printf("DEBUG: Failed to save batch - %v\n", err)
						failed += len(rekaps)
					} else {
						saved += len(rekaps)
					}
					rekaps = rekaps[:0]
					mu.Unlock()
				}

				currentTipeProduksi = newTipe
				fmt.Printf("DEBUG Row %d: Category changed to '%s'\n", i, currentTipeProduksi)
			}
			continue
		}

		if isLikelySummaryRow(row, baseIdx) {
			continue
		}

		if !isValidDataRow(row, baseIdx) {
			continue
		}

		if !hasValidHKO(row, baseIdx) {
			continue
		}

		rekap, err := mapRowRelative(row, baseIdx, tanggal, currentTipeProduksi, afdeling, idMaster)
		if err != nil {
			failed++
			continue
		}

		rekaps = append(rekaps, rekap)

		// Save batch when reaching batch size
		if len(rekaps) >= batchSize {
			mu.Lock()
			if err := saveBatchRekap(db, rekaps); err != nil {
				fmt.Printf("DEBUG: Failed to save batch - %v\n", err)
				failed += len(rekaps)
			} else {
				saved += len(rekaps)
			}
			rekaps = rekaps[:0]
			mu.Unlock()
		}
	}

	// Save remaining records
	if len(rekaps) > 0 {
		if err := saveBatchRekap(db, rekaps); err != nil {
			fmt.Printf("DEBUG: Failed to save final batch - %v\n", err)
			failed += len(rekaps)
		} else {
			saved += len(rekaps)
		}
	}

	fmt.Printf("\nSUMMARY: Saved=%d, Failed=%d, Total Processed=%d\n", saved, failed, saved+failed)
	return saved, failed, nil
}

// ConvertCSVAutoBaseWithFilter: public function to process all CSVs in csv/ folder
func ConvertCSVAutoBaseWithFilter(tanggal time.Time, afdeling string, idMaster uint64) (int, int, []string, error) {
	db := config.GetDB()
	if db == nil {
		return 0, 0, nil, fmt.Errorf("database belum dikonfigurasi (config.GetDB() == nil)")
	}

	path := filepath.Join("csv", "REKAP.csv")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return 0, 0, nil, fmt.Errorf("file REKAP.csv tidak ditemukan di folder csv/")
	}

	var errors []string
	saved, failed, err := processCSVFileAutoBaseWithFilter(db, path, tanggal, afdeling, idMaster)
	if err != nil {
		errors = append(errors, fmt.Sprintf("REKAP.csv: %v", err))
	}

	return saved, failed, errors, nil
}
