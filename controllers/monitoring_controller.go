package controllers

import (
	"app-inputan-ptpn/config"
	"app-inputan-ptpn/models"
	"net/http"
	"time"
)

func GetBakuPenyadapToday(w http.ResponseWriter, r *http.Request) {
	// Ambil tanggal hari ini (format YYYY-MM-DD)
	tanggal := time.Now().Format("2006-01-02")

	// Ambil data penyadap hanya untuk tanggal hari ini
	var penyadap []models.BakuPenyadap
	query := config.DB.Preload("Mandor").Preload("Penyadap").
		Where("DATE(tanggal) = ?", tanggal).
		Order("created_at desc")

	if err := query.Find(&penyadap).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal mengambil data penyadap: " + err.Error(),
		})
		return
	}

	// Kirim response dengan struktur sama seperti GetAllBakuPenyadap
	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Data penyadap untuk tanggal " + tanggal + " berhasil diambil",
		Data:    penyadap,
	})
}
func ServeMonitoringPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/html/monitoring.html")
}
