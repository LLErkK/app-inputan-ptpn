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
func CreateBakuPenyadap(w http.ResponseWriter, r *http.Request) {
	/*
		post baku mandor dan penyadap
		siapa mandornya?
		kemudian nama,basah latek,sheet. basah lump, brcr penyadap
		untuk tanggal tergantung inputan dilakukan tanggal berapa
	*/
}

func DeleteBakuPenyadap(w http.ResponseWriter, r *http.Request) {
	/*manghapus salah satu data penyadap*/
}

func EditBakuPenyadap(w http.ResponseWriter, r *http.Request) {
	/*Mengedit salah satu data terpilih*/
}

func HitungbakuDetail(w http.ResponseWriter, r *http.Request) {
	/*
		setiap kali menambahkan atau menghapus data, data detailnya dihitung lagi
	*/
}
