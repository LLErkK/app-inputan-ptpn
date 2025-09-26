package controllers

import (
	"app-inputan-ptpn/config"
	"app-inputan-ptpn/models"
	"fmt"
	"net/http"
	"time"
)

func ServeRekapPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/html/rekap.html")
}
func GetBakuDetailToday(w http.ResponseWriter, r *http.Request) {
	// Ambil tanggal hari ini (format YYYY-MM-DD)
	today := time.Now().Format("2006-01-02")

	var details []models.BakuDetail
	query := config.DB.
		Where("DATE(tanggal) = DATE(?)", today).
		Order("mandor asc")

	if err := query.Find(&details).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal mengambil detail: " + err.Error(),
		})
		return
	}

	if len(details) == 0 {
		message := "Detail untuk tanggal " + today + " tidak ditemukan"
		respondJSON(w, http.StatusNotFound, APIResponse{
			Success: false,
			Message: message,
		})
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Detail berhasil ditemukan",
		Data:    details,
	})
}
func GetBakuDetailUntilTodayThisMonth(w http.ResponseWriter, r *http.Request) {
	now := time.Now()

	// Awal bulan
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := now

	var details []models.BakuDetail
	if err := config.DB.
		Where("tanggal BETWEEN ? AND ?", startOfMonth, endOfMonth).
		Order("mandor asc, tahun_tanam asc, tipe asc").
		Find(&details).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal mengambil detail: " + err.Error(),
		})
		return
	}

	if len(details) == 0 {
		respondJSON(w, http.StatusNotFound, APIResponse{
			Success: false,
			Message: "Belum ada data baku detail untuk bulan ini hingga " + now.Format("2006-01-02"),
		})
		return
	}

	// Map untuk mengelompokkan data berdasarkan Mandor + TahunTanam + Tipe
	type Rekap struct {
		Mandor     string
		TahunTanam uint
		Tipe       models.TipeProduksi

		JumlahPabrikBasahLatek      float64
		JumlahKebunBasahLatek       float64
		SelisihBasahLatek           float64
		PersentaseSelisihBasahLatek float64

		JumlahSheet float64
		K3Sheet     float64

		JumlahPabrikBasahLump      float64
		JumlahKebunBasahLump       float64
		SelisihBasahLump           float64
		PersentaseSelisihBasahLump float64

		JumlahBrCr float64
		K3BrCr     float64
	}

	rekapMap := make(map[string]*Rekap)

	for _, d := range details {
		key := fmt.Sprintf("%s-%d-%s", d.Mandor, d.TahunTanam, d.Tipe)
		if _, ok := rekapMap[key]; !ok {
			rekapMap[key] = &Rekap{
				Mandor:     d.Mandor,
				TahunTanam: d.TahunTanam,
				Tipe:       d.Tipe,
			}
		}

		r := rekapMap[key]
		r.JumlahPabrikBasahLatek += d.JumlahPabrikBasahLatek
		r.JumlahKebunBasahLatek += d.JumlahKebunBasahLatek
		r.SelisihBasahLatek += d.SelisihBasahLatek
		r.PersentaseSelisihBasahLatek += d.PersentaseSelisihBasahLatek

		r.JumlahSheet += d.JumlahSheet
		r.K3Sheet += d.K3Sheet

		r.JumlahPabrikBasahLump += d.JumlahPabrikBasahLump
		r.JumlahKebunBasahLump += d.JumlahKebunBasahLump
		r.SelisihBasahLump += d.SelisihBasahLump
		r.PersentaseSelisihBasahLump += d.PersentaseSelisihBasahLump

		r.JumlahBrCr += d.JumlahBrCr
		r.K3BrCr += d.K3BrCr
	}

	// Ubah map jadi slice agar mudah dibaca di JSON
	var rekapList []Rekap
	for _, v := range rekapMap {
		rekapList = append(rekapList, *v)
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Rekap baku detail per Mandor, TahunTanam, dan Tipe bulan ini sampai " + now.Format("2006-01-02"),
		Data: map[string]interface{}{
			"start":     startOfMonth.Format("2006-01-02"),
			"end":       endOfMonth.Format("2006-01-02"),
			"totalData": len(details),
			"rekap":     rekapList,
		},
	})
}
