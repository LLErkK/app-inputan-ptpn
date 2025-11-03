package controllers

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"app-inputan-ptpn/config"
	"app-inputan-ptpn/models"

	"gorm.io/gorm"
)

// =====================
// Controller CSV -> DB (deteksi base index + filter summary rows + tipe produksi)
// =====================

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

// detectTipeProduksi: deteksi tipe produksi dari baris kategori
// Mengembalikan tipe produksi yang terdeteksi atau string kosong
func detectTipeProduksi(row []string, baseIdx int) string {
	// Cek di sekitar kolom baseIdx untuk teks tipe produksi
	for i := 0; i < len(row); i++ {
		cell := strings.ToUpper(strings.TrimSpace(row[i]))
		if cell == "" {
			continue
		}

		// Deteksi berbagai variasi tipe produksi
		if strings.Contains(cell, "PRODUKSI BAKU BORONG") {
			return "PRODUKSI BAKU BORONG"
		}
		if strings.Contains(cell, "PRODUKSI BAKU") && !strings.Contains(cell, "BORONG") {
			return "PRODUKSI BAKU"
		}
		if strings.Contains(cell, "PRODUKSI BORONG MINGGU") || strings.Contains(cell, "BORONG MINGGU") {
			return "PRODUKSI BORONG MINGGU"
		}
		if strings.Contains(cell, "PRODUKSI BORONG INTERNAL") || strings.Contains(cell, "BORONG INTERNAL") {
			return "PRODUKSI BORONG INTERNAL"
		}
		if strings.Contains(cell, "PRODUKSI BORONG EKSTERNAL") || strings.Contains(cell, "BORONG EKSTERNAL") {
			return "PRODUKSI BORONG EKSTERNAL"
		}
		if strings.Contains(cell, "PRODUKSI TETES LANJUT") || strings.Contains(cell, "TETES LANJUT") {
			return "PRODUKSI TETES LANJUT"
		}
	}
	return ""
}

// isTipeProduksiRow: cek apakah baris ini adalah baris kategori tipe produksi
func isTipeProduksiRow(row []string, baseIdx int) bool {
	tipe := detectTipeProduksi(row, baseIdx)
	return tipe != ""
}

// isLikelySummaryRow: cek apakah row adalah baris ringkasan
func isLikelySummaryRow(row []string, baseIdx int) bool {
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

// mapRowRelative: mapping urut sesuai struktur model Rekap
// Urutan kolom dari CSV (setelah baseIdx):
// 0: TAHUN TANAM
// 1: NIK
// 2: MANDOR
// 3: HKO HR INI
// 4: HKO S/D HR INI
// 5: BASAH LATEK KEBUN (HR INI)
// 6: BASAH LATEK PABRIK (HR INI)
// 7: BASAH LATEK % (HR INI)
// 8: BASAH LUMP KEBUN (HR INI)
// 9: BASAH LUMP PABRIK (HR INI)
// 10: BASAH LUMP % (HR INI)
// 11: K3 SHEET (HR INI)
// 12: KERING SHEET (HR INI)
// 13: KERING BR.CR (HR INI)
// 14: KERING JUMLAH (HR INI)
// 15: BASAH LATEK KEBUN (S/D HR INI)
// 16: BASAH LATEK PABRIK (S/D HR INI)
// 17: BASAH LATEK % (S/D HR INI)
// 18: BASAH LUMP KEBUN (S/D HR INI)
// 19: BASAH LUMP PABRIK (S/D HR INI)
// 20: BASAH LUMP % (S/D HR INI)
// 21: K3 SHEET (S/D HR INI)
// 22: KERING SHEET (S/D HR INI)
// 23: KERING BR.CR (S/D HR INI)
// 24: KERING JUMLAH (S/D HR INI)
// 25: PRODUKSI PER TAPER HR INI
// 26: PRODUKSI PER TAPER S/D HR INI
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

	// Set tanggal dan tipe produksi
	rekap.Tanggal = tanggal
	rekap.TipeProduksi = tipeProduksi

	// Data identitas (kolom 0-2 dari baseIdx)
	rekap.TahunTanam = getStr(baseIdx + 0)
	rekap.NIK = getStr(baseIdx + 1)
	rekap.Mandor = getStr(baseIdx + 2)

	// HKO (kolom 3-4)
	rekap.HKOHariIni = getInt(baseIdx + 3)
	rekap.HKOSampaiHariIni = getInt(baseIdx + 4)

	// Produksi Hari Ini - Basah Latek (kolom 5-7)
	rekap.HariIniBasahLatekKebun = getFloat(baseIdx + 5)
	rekap.HariIniBasahLatekPabrik = getFloat(baseIdx + 6)
	rekap.HariIniBasahLatekPersen = getFloat(baseIdx + 7)

	// Produksi Hari Ini - Basah Lump (kolom 8-10)
	rekap.HariIniBasahLumpKebun = getFloat(baseIdx + 8)
	rekap.HariIniBasahLumpPabrik = getFloat(baseIdx + 9)
	rekap.HariIniBasahLumpPersen = getFloat(baseIdx + 10)

	// Produksi Hari Ini - Kering (kolom 11-14)
	rekap.HariIniK3Sheet = getFloat(baseIdx + 11)
	rekap.HariIniKeringSheet = getFloat(baseIdx + 12)
	rekap.HariIniKeringBrCr = getFloat(baseIdx + 13)
	rekap.HariIniKeringJumlah = getFloat(baseIdx + 14)

	// Produksi Sampai Hari Ini - Basah Latek (kolom 15-17)
	rekap.SampaiHariIniBasahLatekKebun = getFloat(baseIdx + 15)
	rekap.SampaiHariIniBasahLatekPabrik = getFloat(baseIdx + 16)
	rekap.SampaiHariIniBasahLatekPersen = getFloat(baseIdx + 17)

	// Produksi Sampai Hari Ini - Basah Lump (kolom 18-20)
	rekap.SampaiHariIniBasahLumpKebun = getFloat(baseIdx + 18)
	rekap.SampaiHariIniBasahLumpPabrik = getFloat(baseIdx + 19)
	rekap.SampaiHariIniBasahLumpPersen = getFloat(baseIdx + 20)

	// Produksi Sampai Hari Ini - Kering (kolom 21-24)
	rekap.SampaiHariIniK3Sheet = getFloat(baseIdx + 21)
	rekap.SampaiHariIniKeringSheet = getFloat(baseIdx + 22)
	rekap.SampaiHariIniKeringBrCr = getFloat(baseIdx + 23)
	rekap.SampaiHariIniKeringJumlah = getFloat(baseIdx + 24)

	// Produktivitas Per Taper (kolom 25-26)
	rekap.ProduksiPerTaperHariIni = getFloat(baseIdx + 25)
	rekap.ProduksiPerTaperSampaiHariIni = getFloat(baseIdx + 26)

	rekap.Afdeling = afdeling

	rekap.IdMaster = idMaster

	return rekap, nil
}

// saveRekap: upsert-like save via GORM
func saveRekap(db *gorm.DB, rekap *models.Rekap) error {
	var existing models.Rekap
	res := db.Where("tanggal = ? AND tipe_produksi = ? AND nik = ? AND mandor = ? AND tahun_tanam = ?",
		rekap.Tanggal, rekap.TipeProduksi, rekap.NIK, rekap.Mandor, rekap.TahunTanam).First(&existing)

	if res.Error == nil {
		rekap.ID = existing.ID
		rekap.CreatedAt = existing.CreatedAt
		return db.Save(rekap).Error
	}
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		return res.Error
	}
	return db.Create(rekap).Error
}

// processCSVFileAutoBaseWithFilter: detect baseIdx, filter summary rows, detect tipe produksi, save valid rows
func processCSVFileAutoBaseWithFilter(db *gorm.DB, path string, tanggal time.Time, afdeling string, idMaster uint64) (int, int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, 0, fmt.Errorf("gagal buka file %s: %w", path, err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = -1
	r.LazyQuotes = true

	rows, err := r.ReadAll()
	if err != nil {
		return 0, 0, fmt.Errorf("gagal baca csv %s: %w", path, err)
	}
	if len(rows) == 0 {
		return 0, 0, fmt.Errorf("file kosong: %s", path)
	}

	headerRow, baseIdx := findHeaderRowAndBaseIndex(rows, 30)
	if headerRow == -1 || baseIdx == -1 {
		baseIdx = -1
		for j, c := range rows[0] {
			if strings.Contains(strings.ToUpper(c), "TAHUN TANAM") ||
				strings.Contains(strings.ToUpper(c), "NIK") {
				baseIdx = j
				headerRow = 0
				break
			}
		}
		if baseIdx == -1 {
			return 0, 0, fmt.Errorf("tidak menemukan header 'TAHUN TANAM' atau 'NIK' di %s", path)
		}
	}

	start := headerRow + 1
	saved, failed := 0, 0
	currentTipeProduksi := "PRODUKSI BAKU" // default tipe produksi

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

		// deteksi tipe produksi baru
		if isTipeProduksiRow(row, baseIdx) {
			newTipe := detectTipeProduksi(row, baseIdx)
			if newTipe != "" {
				currentTipeProduksi = newTipe
			}
			continue
		}

		// skip likely summary rows
		if isLikelySummaryRow(row, baseIdx) {
			continue
		}

		// validate minimal data pattern
		if !isValidDataRow(row, baseIdx) {
			continue
		}

		rekap, err := mapRowRelative(row, baseIdx, tanggal, currentTipeProduksi, afdeling, idMaster)
		if err != nil {
			failed++
			continue
		}

		if err := saveRekap(db, rekap); err != nil {
			failed++
			continue
		}
		saved++
	}

	return saved, failed, nil
}

// ConvertCSVAutoBaseWithFilter: public function to process all CSVs in csv/ folder
func ConvertCSVAutoBaseWithFilter(tanggal time.Time, afdeling string, idMaster uint64) (int, int, []string, error) {
	db := config.GetDB()
	if db == nil {
		return 0, 0, nil, fmt.Errorf("database belum dikonfigurasi (config.GetDB() == nil)")
	}

	// Hanya proses file REKAP.csv
	path := filepath.Join("csv", "REKAP.csv")

	// Cek apakah file ada
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
