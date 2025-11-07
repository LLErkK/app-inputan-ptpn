package dev

import (
	"app-inputan-ptpn/config"
	"app-inputan-ptpn/models"
	"encoding/json"
	"net/http"
)

func GetAllProduksi(w http.ResponseWriter, r *http.Request) {
	db := config.GetDB()

	var produksi []models.Produksi
	if err := db.Preload("Master").Find(&produksi).Error; err != nil {
		http.Error(w, "Gagal mengambil data rekap: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(produksi)
}
