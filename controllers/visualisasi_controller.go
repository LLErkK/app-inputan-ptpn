package controllers

import (
	"app-inputan-ptpn/config"
	"app-inputan-ptpn/models"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type BakuSummary struct {
	ID       uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Tanggal  time.Time `json:"tanggal"`
	Afdeling string    `json:"afdeling"`

	JumlahPabrikBasahLatek      float64 `json:"jumlah_pabrik_basah_latek"`
	JumlahKebunBasahLatek       float64 `json:"jumlah_kebun_basah_latek"`
	SelisihBasahLatek           float64 `json:"selisih_basah_latek"`
	PersentaseSelisihBasahLatek float64 `json:"persentase_selisih_basah_latek"`

	JumlahSheet float64 `json:"jumlah_sheet"`
	K3Sheet     float64 `json:"k3_sheet"`

	JumlahPabrikBasahLump      float64 `json:"jumlah_pabrik_basah_lump"`
	JumlahKebunBasahLump       float64 `json:"jumlah_kebun_basah_lump"`
	SelisihBasahLump           float64 `json:"selisih_basah_lump"`
	PersentaseSelisihBasahLump float64 `json:"persentase_selisih_basah_lump"`

	JumlahBrCr float64 `json:"jumlah_br_cr"`
	K3BrCr     float64 `json:"k3_br_cr"`
}

type PenyadapSummary struct {
	ID       uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Tanggal  time.Time `json:"tanggal"`
	IdMandor uint      `json:"id_mandor"`

	BasahLatex float64 `gorm:"default:0" json:"basahLatex"`
	Sheet      float64 `gorm:"default:0" json:"sheet"`
	BasahLump  float64 `gorm:"default:0" json:"basahLump"`
	BrCr       float64 `gorm:"default:0" json:"brCr"`
}

// Helper function untuk mendapatkan range tanggal dari parameter atau default bulan ini
func getDateRange(r *http.Request) (time.Time, time.Time) {
	now := time.Now()

	// Ambil parameter bulan dan tahun dari query string
	bulanStr := r.URL.Query().Get("bulan")
	tahunStr := r.URL.Query().Get("tahun")

	// Default ke bulan dan tahun saat ini
	bulan := int(now.Month())
	tahun := now.Year()

	// Parse bulan jika ada
	if bulanStr != "" {
		if b, err := strconv.Atoi(bulanStr); err == nil && b >= 1 && b <= 12 {
			bulan = b
		}
	}

	// Parse tahun jika ada
	if tahunStr != "" {
		if t, err := strconv.Atoi(tahunStr); err == nil && t > 0 {
			tahun = t
		}
	}

	// Hitung first of month dan end of month
	firstOfMonth := time.Date(tahun, time.Month(bulan), 1, 0, 0, 0, 0, now.Location())
	firstOfNextMonth := time.Date(tahun, time.Month(bulan)+1, 1, 0, 0, 0, 0, now.Location())
	endOfMonth := firstOfNextMonth.AddDate(0, 0, -1)

	return firstOfMonth, endOfMonth
}

func ServeVisualisasiPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/html/visualisasi.html")
}

func GetTotalPerDay(w http.ResponseWriter, r *http.Request) {
	var summaries []BakuSummary
	tipeFilter := r.URL.Query().Get("tipe")

	// Gunakan helper function untuk mendapatkan range tanggal
	firstOfMonth, endOfMonth := getDateRange(r)

	// query dasar
	query := config.DB.
		Model(&models.BakuDetail{}).
		Select(`
            tanggal, afdeling,
            SUM(jumlah_pabrik_basah_latek) as jumlah_pabrik_basah_latek,
            SUM(jumlah_kebun_basah_latek) as jumlah_kebun_basah_latek,
            SUM(selisih_basah_latek) as selisih_basah_latek,
            SUM(persentase_selisih_basah_latek) as persentase_selisih_basah_latek,
            SUM(jumlah_sheet) as jumlah_sheet,
            SUM(k3_sheet) as k3_sheet,
            SUM(jumlah_pabrik_basah_lump) as jumlah_pabrik_basah_lump,
            SUM(jumlah_kebun_basah_lump) as jumlah_kebun_basah_lump,
            SUM(selisih_basah_lump) as selisih_basah_lump,
            SUM(persentase_selisih_basah_lump) as persentase_selisih_basah_lump,
            SUM(jumlah_br_cr) as jumlah_br_cr,
            SUM(k3_br_cr) as k3_br_cr
        `).
		Where("tanggal BETWEEN ? AND ?", firstOfMonth, endOfMonth).
		Group("tanggal, afdeling").
		Order("tanggal asc, afdeling asc")

	// filter tambahan kalau ada parameter tipe
	if tipeFilter != "" {
		query = query.Where("tipe = ?", tipeFilter)
	}

	// eksekusi query
	if err := query.Find(&summaries).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal mengambil data: " + err.Error(),
		})
		return
	}

	// tambahkan ID manual (1,2,3,...)
	for i := range summaries {
		summaries[i].ID = uint(i + 1)
	}

	msg := "Berhasil mengambil data"
	if tipeFilter != "" {
		msg += " bertipe " + tipeFilter
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: msg,
		Data:    summaries,
	})
}

func GetPenyadapPerDay(w http.ResponseWriter, r *http.Request) {
	//parameter id penyadap
	idPenyadap := r.URL.Query().Get("idPenyadap")

	if idPenyadap == "" {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "tidak ada parameter id penyadap",
		})
		return
	}

	var summaries []PenyadapSummary

	// Gunakan helper function untuk mendapatkan range tanggal
	firstOfMonth, endOfMonth := getDateRange(r)

	// query dasar + filter id penyadap
	query := config.DB.
		Model(&models.BakuPenyadap{}).
		Select(`
			tanggal,
			id_baku_mandor as id_mandor,
			SUM(basah_latex) as basah_latex,
			SUM(sheet) as sheet,
			SUM(basah_lump) as basah_lump,
			SUM(br_cr) as br_cr
		`).
		Where("tanggal BETWEEN ? AND ?", firstOfMonth, endOfMonth).
		Where("id_penyadap = ?", idPenyadap).
		Group("tanggal, id_baku_mandor").
		Order("tanggal asc")

	// eksekusi query
	if err := query.Find(&summaries).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal mengambil data: " + err.Error(),
		})
		return
	}

	// tambahkan ID manual (1,2,3,...)
	for i := range summaries {
		summaries[i].ID = uint(i + 1)
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Berhasil mengambil data",
		Data:    summaries,
	})
}

func GetMandorPerDay(w http.ResponseWriter, r *http.Request) {
	//parameter id mandor
	var summaries []BakuSummary
	idMandor := r.URL.Query().Get("idMandor")

	// Gunakan helper function untuk mendapatkan range tanggal
	firstOfMonth, endOfMonth := getDateRange(r)

	// query dasar
	query := config.DB.
		Model(&models.BakuDetail{}).
		Select(`
            tanggal, afdeling,
            SUM(jumlah_pabrik_basah_latek) as jumlah_pabrik_basah_latek,
            SUM(jumlah_kebun_basah_latek) as jumlah_kebun_basah_latek,
            SUM(selisih_basah_latek) as selisih_basah_latek,
            SUM(persentase_selisih_basah_latek) as persentase_selisih_basah_latek,
            SUM(jumlah_sheet) as jumlah_sheet,
            SUM(k3_sheet) as k3_sheet,
            SUM(jumlah_pabrik_basah_lump) as jumlah_pabrik_basah_lump,
            SUM(jumlah_kebun_basah_lump) as jumlah_kebun_basah_lump,
            SUM(selisih_basah_lump) as selisih_basah_lump,
            SUM(persentase_selisih_basah_lump) as persentase_selisih_basah_lump,
            SUM(jumlah_br_cr) as jumlah_br_cr,
            SUM(k3_br_cr) as k3_br_cr
        `).
		Where("tanggal BETWEEN ? AND ? AND id_baku_mandor = ?", firstOfMonth, endOfMonth, idMandor).
		Group("tanggal, afdeling").
		Order("tanggal asc, afdeling asc")

	// eksekusi query
	if err := query.Find(&summaries).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal mengambil data: " + err.Error(),
		})
		return
	}

	// tambahkan ID manual (1,2,3,...)
	for i := range summaries {
		summaries[i].ID = uint(i + 1)
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "berhasil mengambil data",
		Data:    summaries,
	})
}

// GetTotalPerYear - Mendapatkan data total per tahun
func GetTotalPerYear(w http.ResponseWriter, r *http.Request) {
	var summaries []BakuSummary
	tipeFilter := r.URL.Query().Get("tipe")

	// Ambil parameter tahun dari query string, default ke tahun saat ini
	now := time.Now()
	tahunStr := r.URL.Query().Get("tahun")
	tahun := now.Year()

	if tahunStr != "" {
		if t, err := strconv.Atoi(tahunStr); err == nil && t > 0 {
			tahun = t
		}
	}

	// Range 1 Januari - 31 Desember tahun yang dipilih
	startOfYear := time.Date(tahun, 1, 1, 0, 0, 0, 0, now.Location())
	endOfYear := time.Date(tahun, 12, 31, 23, 59, 59, 0, now.Location())

	// query dasar
	query := config.DB.
		Model(&models.BakuDetail{}).
		Select(`
            tanggal, afdeling,
            SUM(jumlah_pabrik_basah_latek) as jumlah_pabrik_basah_latek,
            SUM(jumlah_kebun_basah_latek) as jumlah_kebun_basah_latek,
            SUM(selisih_basah_latek) as selisih_basah_latek,
            SUM(persentase_selisih_basah_latek) as persentase_selisih_basah_latek,
            SUM(jumlah_sheet) as jumlah_sheet,
            SUM(k3_sheet) as k3_sheet,
            SUM(jumlah_pabrik_basah_lump) as jumlah_pabrik_basah_lump,
            SUM(jumlah_kebun_basah_lump) as jumlah_kebun_basah_lump,
            SUM(selisih_basah_lump) as selisih_basah_lump,
            SUM(persentase_selisih_basah_lump) as persentase_selisih_basah_lump,
            SUM(jumlah_br_cr) as jumlah_br_cr,
            SUM(k3_br_cr) as k3_br_cr
        `).
		Where("tanggal BETWEEN ? AND ?", startOfYear, endOfYear).
		Group("tanggal, afdeling").
		Order("tanggal asc, afdeling asc")

	// filter tambahan kalau ada parameter tipe
	if tipeFilter != "" {
		query = query.Where("tipe = ?", tipeFilter)
	}

	// eksekusi query
	if err := query.Find(&summaries).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal mengambil data: " + err.Error(),
		})
		return
	}

	// tambahkan ID manual (1,2,3,...)
	for i := range summaries {
		summaries[i].ID = uint(i + 1)
	}

	msg := "Berhasil mengambil data tahun " + strconv.Itoa(tahun)
	if tipeFilter != "" {
		msg += " bertipe " + tipeFilter
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: msg,
		Data:    summaries,
	})
}

// GetPenyadapPerYear - Mendapatkan data penyadap per tahun
func GetPenyadapPerYear(w http.ResponseWriter, r *http.Request) {
	idPenyadap := r.URL.Query().Get("idPenyadap")

	if idPenyadap == "" {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "tidak ada parameter id penyadap",
		})
		return
	}

	var summaries []PenyadapSummary

	// Ambil parameter tahun dari query string, default ke tahun saat ini
	now := time.Now()
	tahunStr := r.URL.Query().Get("tahun")
	tahun := now.Year()

	if tahunStr != "" {
		if t, err := strconv.Atoi(tahunStr); err == nil && t > 0 {
			tahun = t
		}
	}

	// Range 1 Januari - 31 Desember tahun yang dipilih
	startOfYear := time.Date(tahun, 1, 1, 0, 0, 0, 0, now.Location())
	endOfYear := time.Date(tahun, 12, 31, 23, 59, 59, 0, now.Location())

	// query dasar + filter id penyadap
	query := config.DB.
		Model(&models.BakuPenyadap{}).
		Select(`
			tanggal,
			id_baku_mandor as id_mandor,
			SUM(basah_latex) as basah_latex,
			SUM(sheet) as sheet,
			SUM(basah_lump) as basah_lump,
			SUM(br_cr) as br_cr
		`).
		Where("tanggal BETWEEN ? AND ?", startOfYear, endOfYear).
		Where("id_penyadap = ?", idPenyadap).
		Group("tanggal, id_baku_mandor").
		Order("tanggal asc")

	// eksekusi query
	if err := query.Find(&summaries).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal mengambil data: " + err.Error(),
		})
		return
	}

	// tambahkan ID manual (1,2,3,...)
	for i := range summaries {
		summaries[i].ID = uint(i + 1)
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Berhasil mengambil data tahun " + strconv.Itoa(tahun),
		Data:    summaries,
	})
}

// GetMandorPerYear - Mendapatkan data mandor per tahun
func GetMandorPerYear(w http.ResponseWriter, r *http.Request) {
	var summaries []BakuSummary
	idMandor := r.URL.Query().Get("idMandor")

	if idMandor == "" {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "tidak ada parameter id mandor",
		})
		return
	}

	// Ambil parameter tahun dari query string, default ke tahun saat ini
	now := time.Now()
	tahunStr := r.URL.Query().Get("tahun")
	tahun := now.Year()

	if tahunStr != "" {
		if t, err := strconv.Atoi(tahunStr); err == nil && t > 0 {
			tahun = t
		}
	}

	// Range 1 Januari - 31 Desember tahun yang dipilih
	startOfYear := time.Date(tahun, 1, 1, 0, 0, 0, 0, now.Location())
	endOfYear := time.Date(tahun, 12, 31, 23, 59, 59, 0, now.Location())

	// query dasar
	query := config.DB.
		Model(&models.BakuDetail{}).
		Select(`
            tanggal, afdeling,
            SUM(jumlah_pabrik_basah_latek) as jumlah_pabrik_basah_latek,
            SUM(jumlah_kebun_basah_latek) as jumlah_kebun_basah_latek,
            SUM(selisih_basah_latek) as selisih_basah_latek,
            SUM(persentase_selisih_basah_latek) as persentase_selisih_basah_latek,
            SUM(jumlah_sheet) as jumlah_sheet,
            SUM(k3_sheet) as k3_sheet,
            SUM(jumlah_pabrik_basah_lump) as jumlah_pabrik_basah_lump,
            SUM(jumlah_kebun_basah_lump) as jumlah_kebun_basah_lump,
            SUM(selisih_basah_lump) as selisih_basah_lump,
            SUM(persentase_selisih_basah_lump) as persentase_selisih_basah_lump,
            SUM(jumlah_br_cr) as jumlah_br_cr,
            SUM(k3_br_cr) as k3_br_cr
        `).
		Where("tanggal BETWEEN ? AND ? AND id_baku_mandor = ?", startOfYear, endOfYear, idMandor).
		Group("tanggal, afdeling").
		Order("tanggal asc, afdeling asc")

	// eksekusi query
	if err := query.Find(&summaries).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal mengambil data: " + err.Error(),
		})
		return
	}

	// tambahkan ID manual (1,2,3,...)
	for i := range summaries {
		summaries[i].ID = uint(i + 1)
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Berhasil mengambil data tahun " + strconv.Itoa(tahun),
		Data:    summaries,
	})
}

// GetTotalPerFiveYear - Mendapatkan data total per 5 tahun
func GetTotalPerFiveYear(w http.ResponseWriter, r *http.Request) {
	var summaries []BakuSummary
	tipeFilter := r.URL.Query().Get("tipe")

	now := time.Now()
	tahunStr := r.URL.Query().Get("tahun")
	tahun := now.Year()

	if tahunStr != "" {
		if t, err := strconv.Atoi(tahunStr); err == nil && t > 0 {
			tahun = t
		}
	}

	// Range dari awal tahun (5 tahun lalu) hingga akhir tahun yang dipilih
	startOfYear := time.Date(tahun-5, 1, 1, 0, 0, 0, 0, now.Location())
	endOfYear := time.Date(tahun, 12, 31, 23, 59, 59, 0, now.Location())

	// query dasar
	query := config.DB.
		Model(&models.BakuDetail{}).
		Select(`
            tanggal, afdeling,
            SUM(jumlah_pabrik_basah_latek) as jumlah_pabrik_basah_latek,
            SUM(jumlah_kebun_basah_latek) as jumlah_kebun_basah_latek,
            SUM(selisih_basah_latek) as selisih_basah_latek,
            SUM(persentase_selisih_basah_latek) as persentase_selisih_basah_latek,
            SUM(jumlah_sheet) as jumlah_sheet,
            SUM(k3_sheet) as k3_sheet,
            SUM(jumlah_pabrik_basah_lump) as jumlah_pabrik_basah_lump,
            SUM(jumlah_kebun_basah_lump) as jumlah_kebun_basah_lump,
            SUM(selisih_basah_lump) as selisih_basah_lump,
            SUM(persentase_selisih_basah_lump) as persentase_selisih_basah_lump,
            SUM(jumlah_br_cr) as jumlah_br_cr,
            SUM(k3_br_cr) as k3_br_cr
        `).
		Where("tanggal BETWEEN ? AND ?", startOfYear, endOfYear).
		Group("tanggal, afdeling").
		Order("tanggal asc, afdeling asc")

	// filter tambahan kalau ada parameter tipe
	if tipeFilter != "" {
		query = query.Where("tipe = ?", tipeFilter)
	}

	// eksekusi query
	if err := query.Find(&summaries).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal mengambil data: " + err.Error(),
		})
		return
	}

	// tambahkan ID manual (1,2,3,...)
	for i := range summaries {
		summaries[i].ID = uint(i + 1)
	}

	msg := fmt.Sprintf("Berhasil mengambil data 5 tahun terakhir hingga %d", tahun)
	if tipeFilter != "" {
		msg += fmt.Sprintf(" bertipe %s", tipeFilter)
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: msg,
		Data:    summaries,
	})
}

// GetPenyadapPerFiveYear - Mendapatkan data penyadap per 5 tahun
func GetPenyadapPerFiveYear(w http.ResponseWriter, r *http.Request) {
	idPenyadap := r.URL.Query().Get("idPenyadap")
	if idPenyadap == "" {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "tidak ada parameter id penyadap",
		})
		return
	}

	var summaries []PenyadapSummary

	now := time.Now()
	tahunStr := r.URL.Query().Get("tahun")
	tahun := now.Year()
	if tahunStr != "" {
		if t, err := strconv.Atoi(tahunStr); err == nil && t > 0 {
			tahun = t
		}
	}

	// Range dari awal tahun (5 tahun lalu) hingga akhir tahun yang dipilih
	startOfYear := time.Date(tahun-5, 1, 1, 0, 0, 0, 0, now.Location())
	endOfYear := time.Date(tahun, 12, 31, 23, 59, 59, 0, now.Location())

	// query dasar + filter id penyadap
	query := config.DB.
		Model(&models.BakuPenyadap{}).
		Select(`
			tanggal,
			id_baku_mandor as id_mandor,
			SUM(basah_latex) as basah_latex,
			SUM(sheet) as sheet,
			SUM(basah_lump) as basah_lump,
			SUM(br_cr) as br_cr
		`).
		Where("tanggal BETWEEN ? AND ?", startOfYear, endOfYear).
		Where("id_penyadap = ?", idPenyadap).
		Group("tanggal, id_baku_mandor").
		Order("tanggal asc")

	// eksekusi query
	if err := query.Find(&summaries).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal mengambil data: " + err.Error(),
		})
		return
	}

	// tambahkan ID manual (1,2,3,...)
	for i := range summaries {
		summaries[i].ID = uint(i + 1)
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: fmt.Sprintf("Berhasil mengambil data 5 tahun terakhir hingga %d", tahun),
		Data:    summaries,
	})
}

// GetMandorPerFiveYear - Mendapatkan data mandor per 5 tahun
func GetMandorPerFiveYear(w http.ResponseWriter, r *http.Request) {
	idMandor := r.URL.Query().Get("idMandor")
	if idMandor == "" {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "tidak ada parameter id mandor",
		})
		return
	}

	var summaries []BakuSummary

	now := time.Now()
	tahunStr := r.URL.Query().Get("tahun")
	tahun := now.Year()
	if tahunStr != "" {
		if t, err := strconv.Atoi(tahunStr); err == nil && t > 0 {
			tahun = t
		}
	}

	// Range dari awal tahun (5 tahun lalu) hingga akhir tahun yang dipilih
	startOfYear := time.Date(tahun-5, 1, 1, 0, 0, 0, 0, now.Location())
	endOfYear := time.Date(tahun, 12, 31, 23, 59, 59, 0, now.Location())

	// query dasar
	query := config.DB.
		Model(&models.BakuDetail{}).
		Select(`
            tanggal, afdeling,
            SUM(jumlah_pabrik_basah_latek) as jumlah_pabrik_basah_latek,
            SUM(jumlah_kebun_basah_latek) as jumlah_kebun_basah_latek,
            SUM(selisih_basah_latek) as selisih_basah_latek,
            SUM(persentase_selisih_basah_latek) as persentase_selisih_basah_latek,
            SUM(jumlah_sheet) as jumlah_sheet,
            SUM(k3_sheet) as k3_sheet,
            SUM(jumlah_pabrik_basah_lump) as jumlah_pabrik_basah_lump,
            SUM(jumlah_kebun_basah_lump) as jumlah_kebun_basah_lump,
            SUM(selisih_basah_lump) as selisih_basah_lump,
            SUM(persentase_selisih_basah_lump) as persentase_selisih_basah_lump,
            SUM(jumlah_br_cr) as jumlah_br_cr,
            SUM(k3_br_cr) as k3_br_cr
        `).
		Where("tanggal BETWEEN ? AND ? AND id_baku_mandor = ?", startOfYear, endOfYear, idMandor).
		Group("tanggal, afdeling").
		Order("tanggal asc, afdeling asc")

	// eksekusi query
	if err := query.Find(&summaries).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal mengambil data: " + err.Error(),
		})
		return
	}

	// tambahkan ID manual (1,2,3,...)
	for i := range summaries {
		summaries[i].ID = uint(i + 1)
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: fmt.Sprintf("Berhasil mengambil data 5 tahun terakhir hingga %d", tahun),
		Data:    summaries,
	})
}
