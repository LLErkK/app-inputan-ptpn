package dev

import (
	"app-inputan-ptpn/config"
	"app-inputan-ptpn/models"
	"encoding/json"
	"net/http"
)

func GetAllRekap(w http.ResponseWriter, r *http.Request) {
	db := config.GetDB()

	var rekaps []models.Rekap

	// Ambil semua data rekap dari database
	if err := db.Preload("Master").Find(&rekaps).Error; err != nil {
		http.Error(w, "Gagal mengambil data rekap: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set header response agar berupa JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rekaps)
}
