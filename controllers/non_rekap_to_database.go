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

	"app-inputan-ptpn/config"
	"app-inputan-ptpn/models"

	"gorm.io/gorm"
)

// =====================
// Controller CSV Tanggal -> DB (format berbeda, tipe produksi dari nama file)
// =====================

// findRowContaining: cari baris yang mengandung kata kunci tertentu
func findRowContaining(rows [][]string, keyword string, maxScan int) int {
	limit := maxScan
	if len(rows) < limit {
		limit = len(rows)
	}
	kw := strings.ToLower(strings.TrimSpace(keyword))
	for i := 0; i < limit; i++ {
		for _, cell := range rows[i] {
			if strings.Contains(strings.ToLower(strings.TrimSpace(cell)), kw) {
				return i
			}
		}
	}
	return -1
}

// extractAvailableDates: ekstrak tanggal yang tersedia dari baris nomor tanggal
func extractAvailableDates(row []string) []int {
	dates := make(map[int]bool)
	re := regexp.MustCompile(`^\d+`)
	for _, val := range row {
		cleaned := strings.TrimSpace(val)
		match := re.FindString(cleaned)
		if match != "" {
			if num, err := strconv.Atoi(match); err == nil {
				if num >= 1 && num <= 31 {
					dates[num] = true
				}
			}
		}
	}

	// Convert map to sorted slice
	result := make([]int, 0, len(dates))
	for date := range dates {
		result = append(result, date)
	}
	// Simple sort
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i] > result[j] {
				result[i], result[j] = result[j], result[i]
			}
		}
	}
	return result
}

// findFirstColumnIndexForDate: cari indeks kolom PERTAMA untuk tanggal yang dipilih
func findFirstColumnIndexForDate(row []string, targetDate int) int {
	targetStr := strconv.Itoa(targetDate)
	re := regexp.MustCompile(`^\d+`)
	for colIdx, val := range row {
		cleaned := strings.TrimSpace(val)
		match := re.FindString(cleaned)
		if match == targetStr {
			return colIdx
		}
		// kadang ada header seperti "01/09" atau "1"
		if strings.Contains(strings.ReplaceAll(strings.ToLower(cleaned), " ", ""), targetStr) {
			return colIdx
		}
	}
	return -1
}

// extractTipeProduksiFromFilename: ekstrak tipe produksi dari nama file
// Contoh: "Pantauan_Produksi_Afd_Setro_27-10-2025_Baku.csv" -> "PRODUKSI BAKU"
func extractTipeProduksiFromFilename(filename string) string {
	filename = strings.ToUpper(filename)

	if strings.Contains(filename, "BAKU BORONG") {
		return "PRODUKSI BAKU BORONG"
	}
	if strings.Contains(filename, "BAKU") && !strings.Contains(filename, "BORONG") {
		return "PRODUKSI BAKU"
	}
	if strings.Contains(filename, "BORONG MINGGU") {
		return "PRODUKSI BORONG MINGGU"
	}
	if strings.Contains(filename, "BORONG INTERNAL") {
		return "PRODUKSI BORONG INTERNAL"
	}
	if strings.Contains(filename, "BORONG EKSTERNAL") {
		return "PRODUKSI BORONG EKSTERNAL"
	}
	if strings.Contains(filename, "TETES LANJUT") {
		return "PRODUKSI TETES LANJUT"
	}

	// Default
	return "PRODUKSI BAKU"
}

// saveProduksi: upsert-like save via GORM untuk model Produksi
func saveProduksi(db *gorm.DB, produksi *models.Produksi) error {
	var existing models.Produksi
	res := db.Where("tanggal = ? AND tipe_produksi = ? AND nik = ? AND mandor = ? AND tahun_tanam = ?",
		produksi.Tanggal, produksi.TipeProduksi, produksi.NIK, produksi.Mandor, produksi.TahunTanam).First(&existing)

	if res.Error == nil {
		produksi.ID = existing.ID
		produksi.CreatedAt = existing.CreatedAt
		return db.Save(produksi).Error
	}
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		return res.Error
	}
	return db.Create(produksi).Error
}

func isIrrelevantRow(mandor string) bool {
	irrelevantKeywords := []string{"jumlah", "pabrik", "k3", "selisih", "%", "total", "rata", "rekap", "jumlah keseluruhan"}
	mandorLower := strings.ToLower(strings.TrimSpace(mandor))

	for _, keyword := range irrelevantKeywords {
		if strings.Contains(mandorLower, keyword) {
			return true
		}
	}
	return false
}

// parseNumber: parse angka dengan format Indonesia (1.234,56) atau internasional (1,234.56)
func parseNumberNoRekap(s string) (float64, error) {
	orig := strings.TrimSpace(s)
	if orig == "" {
		return 0, fmt.Errorf("empty")
	}

	// Remove surrounding quotes and non-breaking spaces
	orig = strings.ReplaceAll(orig, "\"", "")
	orig = strings.ReplaceAll(orig, "\u00A0", "")
	orig = strings.TrimSpace(orig)

	// Handle parentheses as negative numbers: (1.234,56) => -1234.56
	neg := false
	if strings.HasPrefix(orig, "(") && strings.HasSuffix(orig, ")") {
		neg = true
		orig = strings.TrimPrefix(strings.TrimSuffix(orig, ")"), "(")
	}

	// Replace common dashes or em dashes with minus
	orig = strings.ReplaceAll(orig, "—", "-")
	orig = strings.ReplaceAll(orig, "–", "-")

	// Remove percentage sign or other trailing non-numeric chars
	orig = strings.TrimRightFunc(orig, func(r rune) bool {
		return !(r >= '0' && r <= '9') && r != '.' && r != ',' && r != '-' && r != '+'
	})

	// If contains both '.' and ',' then assume '.' thousands separator and ',' decimal (Indonesian)
	if strings.Contains(orig, ".") && strings.Contains(orig, ",") {
		orig = strings.ReplaceAll(orig, ".", "")  // remove thousands
		orig = strings.ReplaceAll(orig, ",", ".") // convert decimal separator
	} else if strings.Count(orig, ",") > 0 && !strings.Contains(orig, ".") {
		// if only comma present, assume comma is decimal separator
		orig = strings.ReplaceAll(orig, ",", ".")
	} else {
		//
	}

	// Remove any remaining spaces
	orig = strings.ReplaceAll(orig, " ", "")

	// Attempt parse
	val, err := strconv.ParseFloat(orig, 64)
	if err != nil {
		return 0, err
	}
	if neg {
		val = -val
	}
	return val, nil
}

// mapRowTanggalFormat: mapping row dengan format tanggal ke model Produksi
// Base columns: No, Tahun Tanam, Mandor, NIK, Nama Penyadap
// Production columns (4 kolom berurutan dari firstColIdx): Basah Latek, Sheet, Basah Lump, Br.Cr
func mapRowTanggalFormat(row []string, baseColIndices map[string]int, firstColIdx int, tanggal time.Time, tipeProduksi string, afdeling string, idMaster uint64) (*models.Produksi, error) {
	produksi := &models.Produksi{}

	getStr := func(colName string) string {
		if idx, ok := baseColIndices[colName]; ok && idx >= 0 && idx < len(row) {
			return strings.TrimSpace(strings.ReplaceAll(row[idx], "\"", ""))
		}
		return ""
	}

	getFloat := func(idx int) float64 {
		if idx < 0 || idx >= len(row) {
			return 0
		}
		v := strings.TrimSpace(strings.ReplaceAll(row[idx], "\"", ""))
		if v == "" || v == "-" || v == "—" || v == "—" {
			return 0
		}
		if n, err := parseNumberNoRekap(v); err == nil {
			return n
		}
		return 0
	}

	// Set tanggal dan tipe produksi
	produksi.Tanggal = tanggal
	produksi.TipeProduksi = tipeProduksi

	// Data identitas
	produksi.TahunTanam = getStr("Tahun Tanam")
	produksi.NIK = getStr("NIK")
	produksi.Mandor = getStr("Mandor")
	produksi.NamaPenyadap = getStr("Nama Penyadap")

	// Production columns (4 kolom berurutan)
	if firstColIdx >= 0 {
		produksi.BasahLatek = getFloat(firstColIdx + 0) // Basah Latek
		produksi.Sheet = getFloat(firstColIdx + 1)      // Sheet
		produksi.BasahLump = getFloat(firstColIdx + 2)  // Basah Lump
		produksi.BrCr = getFloat(firstColIdx + 3)       // Br.Cr
	}
	produksi.TotalProduksi = produksi.Sheet + produksi.BasahLump

	produksi.Afdeling = afdeling
	produksi.IdMaster = idMaster

	return produksi, nil
}

// processCSVTanggalFormat: proses CSV dengan format tanggal
func processCSVTanggalFormat(db *gorm.DB, path string, targetDate int, tipeProduksi string, afdeling string, idMaster uint64) (int, int, error) {
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

	// 1. Cari baris "Tanggal"
	tanggalRowIdx := findRowContaining(rows, "tanggal", 15)
	if tanggalRowIdx == -1 {
		// kadang header tanggal tidak eksplisit; coba cari baris yang mengandung angka 1..31 dengan beberapa kolom
		tanggalRowIdx = -1
		for i := 0; i < min(15, len(rows)); i++ {
			// hitung berapa cell yang mendekati angka tanggal
			count := 0
			for _, cell := range rows[i] {
				if matched, _ := regexp.MatchString(`^\s*\d{1,2}\s*$`, strings.TrimSpace(cell)); matched {
					count++
				}
			}
			if count >= 3 { // asumsi baris nomor tanggal biasanya punya beberapa angka
				tanggalRowIdx = i
				break
			}
		}
		if tanggalRowIdx == -1 {
			return 0, 0, fmt.Errorf("tidak menemukan baris 'Tanggal' di %s", path)
		}
	}

	// 2. Baris nomor tanggal (baris setelah "Tanggal")
	numberRowIdx := tanggalRowIdx + 1
	if numberRowIdx >= len(rows) {
		return 0, 0, fmt.Errorf("tidak ada baris nomor tanggal di %s", path)
	}
	numberRow := rows[numberRowIdx]

	// 3. Ekstrak tanggal yang tersedia
	availableDates := extractAvailableDates(numberRow)
	if len(availableDates) == 0 {
		// coba juga ambil dari numberRow sendiri tanpa regex ketat
		for _, cell := range numberRow {
			clean := strings.TrimSpace(cell)
			if n, err := strconv.Atoi(strings.TrimLeft(clean, "0")); err == nil && n >= 1 && n <= 31 {
				availableDates = append(availableDates, n)
			}
		}
		if len(availableDates) == 0 {
			return 0, 0, fmt.Errorf("tidak ada tanggal valid di %s", path)
		}
	}

	// 4. Validasi target date tersedia
	dateFound := false
	for _, d := range availableDates {
		if d == targetDate {
			dateFound = true
			break
		}
	}
	if !dateFound {
		return 0, 0, fmt.Errorf("tanggal %d tidak ditemukan di %s (tersedia: %v)", targetDate, path, availableDates)
	}

	// 5. Cari indeks PERTAMA dari tanggal yang dipilih
	firstColIdx := findFirstColumnIndexForDate(numberRow, targetDate)
	if firstColIdx == -1 {
		return 0, 0, fmt.Errorf("tidak dapat menemukan kolom untuk tanggal %d di %s", targetDate, path)
	}

	// 6. Cari baris "Basah Latek" di area dekat numberRow
	subRowsEnd := min(numberRowIdx+6, len(rows))
	basahLatekRowRel := findRowContaining(rows[numberRowIdx:subRowsEnd], "basah latek", 6)
	basahLatekRow := -1
	if basahLatekRowRel != -1 {
		basahLatekRow = numberRowIdx + basahLatekRowRel
	}

	// 7. Cari baris "Tahun Tanam" untuk header (lebih toleran ke variasi)
	headerStartIdx := findRowContaining(rows, "tahun tanam", len(rows))
	if headerStartIdx == -1 {
		// coba cari "tahun" atau "thn tanam" atau baris yang mengandung "nik" & "mandor"
		for i := 0; i < len(rows); i++ {
			rowLower := strings.ToLower(strings.Join(rows[i], " "))
			if strings.Contains(rowLower, "nik") && (strings.Contains(rowLower, "mandor") || strings.Contains(rowLower, "nama")) {
				headerStartIdx = i
				break
			}
			if strings.Contains(rowLower, "tahun") && strings.Contains(rowLower, "tanam") {
				headerStartIdx = i
				break
			}
		}
		if headerStartIdx == -1 {
			return 0, 0, fmt.Errorf("tidak menemukan header 'Tahun Tanam' di %s", path)
		}
	}

	// 8. Tentukan baris awal data
	dataStartIdx := headerStartIdx + 1
	if basahLatekRow > headerStartIdx {
		dataStartIdx = basahLatekRow + 1
	}

	// 9. Build base column indices
	headerRow := rows[headerStartIdx]
	baseColIndices := make(map[string]int)
	for idx, cell := range headerRow {
		cellClean := strings.TrimSpace(cell)
		cellLower := strings.ToLower(cellClean)

		switch {
		case strings.Contains(cellLower, "tahun tanam") || cellLower == "tahun" || strings.Contains(cellLower, "thn tanam"):
			baseColIndices["Tahun Tanam"] = idx
		case cellLower == "nik" || strings.Contains(cellLower, "no nik"):
			baseColIndices["NIK"] = idx
		case strings.Contains(cellLower, "mandor") || strings.Contains(cellLower, "nama mandor"):
			baseColIndices["Mandor"] = idx
		case strings.Contains(cellLower, "nama penyadap") || strings.Contains(cellLower, "nama") || strings.Contains(cellLower, "penyadap"):
			baseColIndices["Nama Penyadap"] = idx
		}
	}

	// Jika ada missing base columns, coba fallback berdasarkan posisi relatif
	if _, ok := baseColIndices["Tahun Tanam"]; !ok && len(headerRow) >= 1 {
		// asumsikan kolom 1 atau 2
		if len(headerRow) > 1 {
			baseColIndices["Tahun Tanam"] = 1
		} else {
			baseColIndices["Tahun Tanam"] = 0
		}
	}
	if _, ok := baseColIndices["NIK"]; !ok {
		// cari kolom yang berisi angka panjang (NIK)
		for idx, cell := range headerRow {
			if strings.Contains(strings.ToLower(cell), "nik") {
				baseColIndices["NIK"] = idx
				break
			}
		}
	}
	// kalau masih belum ada, biarkan kosong (akan dicek saat parsing rows)

	// 10. Parse tanggal dari nama file atau gunakan target date
	// Format: Pantauan_Produksi_Afd_Setro_27-10-2025_Baku.csv
	var tanggal time.Time
	filename := filepath.Base(path)
	datePattern := regexp.MustCompile(`(\d{2})-(\d{2})-(\d{4})`)
	if match := datePattern.FindStringSubmatch(filename); match != nil {
		day, _ := strconv.Atoi(match[1])
		month, _ := strconv.Atoi(match[2])
		year, _ := strconv.Atoi(match[3])
		tanggal = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	} else {
		// Fallback: gunakan target date dengan bulan/tahun saat ini (lokal)
		now := time.Now()
		tanggal = time.Date(now.Year(), now.Month(), targetDate, 0, 0, 0, 0, time.Local)
	}

	// 11. Process data rows
	saved, failed := 0, 0
	var lastTahunTanam, lastMandor string

	for i := dataStartIdx; i < len(rows); i++ {
		row := rows[i]

		// Skip empty rows
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

		// Forward fill Tahun Tanam dan Mandor
		tahunTanam := ""
		if idx, ok := baseColIndices["Tahun Tanam"]; ok && idx < len(row) {
			tahunTanam = strings.TrimSpace(row[idx])
		}
		if tahunTanam != "" {
			lastTahunTanam = tahunTanam
		} else if idx, ok := baseColIndices["Tahun Tanam"]; ok && idx < len(row) {
			// isi kembali pada row agar mapRowTanggalFormat mendapat value
			row[idx] = lastTahunTanam
		}

		mandor := ""
		if idx, ok := baseColIndices["Mandor"]; ok && idx < len(row) {
			mandor = strings.TrimSpace(row[idx])
		}
		if mandor != "" {
			lastMandor = mandor
		} else if idx, ok := baseColIndices["Mandor"]; ok && idx < len(row) {
			row[idx] = lastMandor
		}

		// Skip irrelevant rows
		if mandor == "" && lastMandor != "" {
			mandor = lastMandor
		}
		if isIrrelevantRow(mandor) {
			continue
		}

		// Check NIK atau Nama Penyadap harus ada
		nik := ""
		namaPenyadap := ""
		if idx, ok := baseColIndices["NIK"]; ok && idx < len(row) {
			nik = strings.TrimSpace(row[idx])
		} else {
			// coba cari kolom NIK heuristik: ada angka panjang
			for ci, cv := range row {
				if len(strings.TrimSpace(cv)) >= 9 && len(strings.TrimSpace(cv)) <= 20 { // heuristik NIK
					// jangan overwrite kalau sudah ada
					if nik == "" {
						nik = strings.TrimSpace(cv)
						baseColIndices["NIK"] = ci
					}
				}
			}
		}
		if idx, ok := baseColIndices["Nama Penyadap"]; ok && idx < len(row) {
			namaPenyadap = strings.TrimSpace(row[idx])
		}

		if nik == "" && namaPenyadap == "" {
			continue
		}

		// Map row to Produksi
		produksi, err := mapRowTanggalFormat(row, baseColIndices, firstColIdx, tanggal, tipeProduksi, afdeling, idMaster)
		if err != nil {
			failed++
			continue
		}

		// Validate minimal data
		if produksi.TahunTanam == "" || produksi.NIK == "" {
			// skip jika data identitas tidak lengkap
			continue
		}
		// Skip jika semua nilai produksi = 0
		if produksi.BasahLatek == 0 && produksi.Sheet == 0 && produksi.BasahLump == 0 && produksi.BrCr == 0 {
			continue
		}

		// Save to database
		if err := saveProduksi(db, produksi); err != nil {
			failed++
			continue
		}
		saved++
	}

	return saved, failed, nil
}

// ConvertCSVTanggalFormat: public function to process all CSVs except REKAP.csv
func ConvertCSVTanggalFormat(targetDate int, afdeling string, idMaster uint64) (int, int, []string, error) {
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

		filename := fi.Name()

		// Skip bukan CSV
		if !strings.HasSuffix(strings.ToLower(filename), ".csv") {
			continue
		}

		// Skip REKAP.csv
		if strings.ToUpper(filename) == "REKAP.CSV" {
			continue
		}

		// Extract tipe produksi dari nama file
		tipeProduksi := extractTipeProduksiFromFilename(filename)

		path := filepath.Join("csv", filename)
		saved, failed, err := processCSVTanggalFormat(db, path, targetDate, tipeProduksi, afdeling, idMaster)
		totalSaved += saved
		totalFailed += failed

		if err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", filename, err))
		}
	}

	if totalSaved == 0 && totalFailed == 0 {
		return 0, 0, nil, fmt.Errorf("tidak ada file CSV yang diproses (selain REKAP.csv)")
	}

	return totalSaved, totalFailed, errors, nil
}

// min helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
