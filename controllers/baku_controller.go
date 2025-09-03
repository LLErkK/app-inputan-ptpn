package controllers

import (
	"app-inputan-ptpn/config"
	"app-inputan-ptpn/models"
	"encoding/json"
	"net/http"
)

type ResponseBaku struct {
	ProduksiBaku       []models.ProduksiBaku       `json:"produksi_baku"`
	ProduksiBakuDetail []models.ProduksibakuDetail `json:"produksi_baku_detail"`
}

func GetAllBaku(w http.ResponseWriter, r *http.Request) {
	var produksiBaku []models.ProduksiBaku
	var produksiBakuDetail []models.ProduksibakuDetail
	config.DB.Order("created_at desc").Find(&produksiBaku)
	config.DB.Order("created_at desc").Find(&produksiBakuDetail)

	response := ResponseBaku{
		ProduksiBaku:       produksiBaku,
		ProduksiBakuDetail: produksiBakuDetail,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
