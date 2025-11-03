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
// Controller CSV -> DB (deteksi base index + filter summary rows)
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

	// keep percent as plain number (12.5% -> 12.5). If you want fraction, divide by 100 here.
	_ = hasPercent

	return num, nil
}

// findHeaderRowAndBaseIndex: cari row & kolom yang mengandung "TAHUN TANAM" (atau "NIK")
// return (headerRowIndex, baseColIndex) or (-1,-1) kalau tidak ditemukan
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
			// normalize common variations
			c := strings.ReplaceAll(clean, "_", " ")
			c = strings.ReplaceAll(c, "\"", "")
			c = strings.TrimSpace(c)
			if strings.Contains(c, "TAHUN TANAM") || c == "TAHUN" || c == "TAHUN_TANAM" {
				return i, j
			}
			if c == "NIK" {
				// jika hanya NIK ditemukan, assume tahun tanam ada di kolom sebelumnya
				return i, j - 1
			}
		}
	}
	return -1, -1
}

// isLikelySummaryRow: cek apakah row pada posisi baseIdx+1 (nik) berisi kata ringkasan seperti "JUMLAH", "SELISIH", "%", "K3", dll.
func isLikelySummaryRow(row []string, baseIdx int) bool {
	checkIdx := baseIdx + 1
	if checkIdx < 0 || checkIdx >= len(row) {
		return false
	}
	cell := strings.ToUpper(strings.TrimSpace(row[checkIdx]))
	if cell == "" {
		// kadang ringkasan berada di kolom tahun_tanam, cek juga
		c0 := strings.ToUpper(strings.TrimSpace(row[baseIdx]))
		if strings.Contains(c0, "JUMLAH") || strings.Contains(c0, "SELISIH") || strings.Contains(c0, "%") || strings.Contains(c0, "K3") {
			return true
		}
		return false
	}
	// kata-kata umum summary
	summaryKeys := []string{"JUMLAH", "SELISIH", "TOTAL", "K3", "%", "JUMLAH PABRIK", "JUMLAH KEBUN", "RATA", "RATA-RATA"}
	for _, k := range summaryKeys {
		if strings.Contains(cell, k) {
			return true
		}
	}
	// jika cell cuma berupa kata non-numeric dan pendek seperti "K3", treat summary
	if regexp.MustCompile(`^[A-Z ]+$`).MatchString(cell) && len(cell) <= 6 && !regexp.MustCompile(`[0-9]`).MatchString(cell) {
		// probably label row like "K3", "TOTAL", "JUMLAH"
		return true
	}
	return false
}

// isValidDataRow: cek minimal validitas baris data
// - tahun_tanam harus angka (4 digit antara 1900..2100)
// - nik harus mengandung digit minimal length (>=5) atau numeric
func isValidDataRow(row []string, baseIdx int) bool {
	// check tahun
	tyIdx := baseIdx
	if tyIdx < 0 || tyIdx >= len(row) {
		return false
	}
	ty := strings.TrimSpace(row[tyIdx])
	if ty == "" {
		return false
	}
	// some headers have "TAHUN TANAM" text; ensure ty is numeric year
	tyDigits := regexp.MustCompile(`\d{4}`)
	if !tyDigits.MatchString(ty) {
		// maybe ty includes extra text but contains 4-digit year -> accept
		return false
	}
	// parse year
	yearStr := tyDigits.FindString(ty)
	if y, err := strconv.Atoi(yearStr); err != nil || y < 1900 || y > 2100 {
		return false
	}

	// check NIK
	nikIdx := baseIdx + 1
	if nikIdx < 0 || nikIdx >= len(row) {
		return false
	}
	nik := strings.TrimSpace(row[nikIdx])
	if nik == "" {
		return false
	}
	// clean nik: remove dots/spaces
	nikClean := strings.ReplaceAll(nik, ".", "")
	nikClean = strings.ReplaceAll(nikClean, " ", "")
	// if contains letters only like "JUMLAH PABRIK" -> invalid
	if !regexp.MustCompile(`\d`).MatchString(nikClean) {
		return false
	}
	// require at least 4-5 digits to be likely NIK/id
	digits := regexp.MustCompile(`\d+`).FindString(nikClean)
	if len(digits) < 4 {
		// but allow cases where NIK is actually name and the id is next column (rare)
		// conservatively mark invalid
		return false
	}

	return true
}

// mapRowRelative: mapping relatif terhadap baseIdx.
// expects: tahun = row[baseIdx], nik=row[baseIdx+1], mandor=row[baseIdx+2], numbers start at baseIdx+3
func mapRowRelative(row []string, baseIdx int, tanggal time.Time) (*models.Rekap, error) {
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
	rekap.TahunTanam = getStr(baseIdx)
	rekap.NIK = getStr(baseIdx + 1)
	rekap.Mandor = getStr(baseIdx + 2)

	// numeric mapping starting at baseIdx+3 (adjust based on CSV)
	k := baseIdx + 3

	rekap.HKOHariIni = getInt(k + 0)
	rekap.HKOSampaiHariIni = getInt(k + 1)

	rekap.HariIniBasahLatekKebun = getFloat(k + 2)
	rekap.HariIniBasahLatekPabrik = getFloat(k + 3)
	rekap.HariIniBasahLatekPersen = getFloat(k + 4)

	rekap.HariIniBasahLumpKebun = getFloat(k + 5)
	rekap.HariIniBasahLumpPabrik = getFloat(k + 6)
	rekap.HariIniBasahLumpPersen = getFloat(k + 7)

	rekap.HariIniK3Sheet = getFloat(k + 8)
	rekap.HariIniKeringSheet = getFloat(k + 9)
	rekap.HariIniKeringBrCr = getFloat(k + 10)
	rekap.HariIniKeringJumlah = getFloat(k + 11)

	rekap.SampaiHariIniBasahLatekKebun = getFloat(k + 12)
	rekap.SampaiHariIniBasahLatekPabrik = getFloat(k + 13)
	rekap.SampaiHariIniBasahLatekPersen = getFloat(k + 14)

	rekap.SampaiHariIniBasahLumpKebun = getFloat(k + 15)
	rekap.SampaiHariIniBasahLumpPabrik = getFloat(k + 16)
	rekap.SampaiHariIniBasahLumpPersen = getFloat(k + 17)

	rekap.SampaiHariIniK3Sheet = getFloat(k + 18)
	rekap.SampaiHariIniKeringSheet = getFloat(k + 19)
	rekap.SampaiHariIniKeringBrCr = getFloat(k + 20)
	rekap.SampaiHariIniKeringJumlah = getFloat(k + 21)

	rekap.ProduksiPerTaperHariIni = getFloat(k + 22)
	rekap.ProduksiPerTaperSampaiHariIni = getFloat(k + 23)

	return rekap, nil
}

// saveRekap: upsert-like save via GORM
func saveRekap(db *gorm.DB, rekap *models.Rekap) error {
	var existing models.Rekap
	res := db.Where("tanggal = ? AND nik = ? AND mandor = ? AND tahun_tanam = ?",
		rekap.Tanggal, rekap.NIK, rekap.Mandor, rekap.TahunTanam).First(&existing)

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

// processCSVFileAutoBaseWithFilter: detect baseIdx, filter summary rows, save valid rows
func processCSVFileAutoBaseWithFilter(db *gorm.DB, path string, tanggal time.Time) (int, int, error) {
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
		// fallback: search first row directly
		baseIdx = -1
		for j, c := range rows[0] {
			if strings.Contains(strings.ToUpper(c), "TAHUN TANAM") || strings.Contains(strings.ToUpper(c), "NIK") {
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

		// skip likely summary rows
		if isLikelySummaryRow(row, baseIdx) {
			continue
		}

		// validate minimal data pattern
		if !isValidDataRow(row, baseIdx) {
			// not a valid data row -> skip (counts as skipped, not fail)
			continue
		}

		rekap, err := mapRowRelative(row, baseIdx, tanggal)
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
func ConvertCSVAutoBaseWithFilter(tanggal time.Time) (int, int, []string, error) {
	files, err := os.ReadDir("csv")
	if err != nil {
		return 0, 0, nil, fmt.Errorf("gagal baca folder csv: %w", err)
	}

	db := config.GetDB()
	if db == nil {
		return 0, 0, nil, fmt.Errorf("database belum dikonfigurasi (config.GetDB() == nil)")
	}

	totalSaved := 0
	totalFailed := 0
	var errors []string

	for _, fi := range files {
		if fi.IsDir() {
			continue
		}
		if !strings.HasSuffix(strings.ToLower(fi.Name()), ".csv") {
			continue
		}
		path := filepath.Join("csv", fi.Name())
		saved, failed, err := processCSVFileAutoBaseWithFilter(db, path, tanggal)
		totalSaved += saved
		totalFailed += failed
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", fi.Name(), err))
		}
	}

	return totalSaved, totalFailed, errors, nil
}
