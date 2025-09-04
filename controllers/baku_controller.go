package controllers

import (
	"app-inputan-ptpn/config"
	"app-inputan-ptpn/models"
	"encoding/json"
	"net/http"
)

func GetAllBakuPenyadap(w http.ResponseWriter, r *http.Request) {
	var bakuPenyadap []models.BakuPenyadap
	config.DB.Order("created_at desc").Find(&bakuPenyadap)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bakuPenyadap)
}
func GetAllBakuMandor(w http.ResponseWriter, r *http.Request) {
	var bakuMandor []models.BakuMandor
	config.DB.Order("created_at desc").Find(&bakuMandor)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bakuMandor)
}
