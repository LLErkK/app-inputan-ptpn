package controllers

import (
	"app-inputan-ptpn/config"
	"app-inputan-ptpn/models"
	"net/http"
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

func ServeVisualisasiPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/html/visualisasi.html")
}
func GetTotalPerDay(w http.ResponseWriter, r *http.Request) {
	var summaries []BakuSummary
	tipeFilter := r.URL.Query().Get("tipe")

	now := time.Now()
	firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	firstOfNextMonth := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())
	endOfMonth := firstOfNextMonth.AddDate(0, 0, -1)

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

func GetTotalPerDayByType(w http.ResponseWriter, r *http.Request) {

}
func GetPenyadapPerDay() {
	//parameter range tanggal
	//parameter id penyadap
}

func GetMandorPerDay() {
	//parameter range tanggal
	//parameter id mandor
}
