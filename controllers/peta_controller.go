package controllers

import (
	"app-inputan-ptpn/config"
	"app-inputan-ptpn/models"
	"encoding/json"
	"net/http"
)

func GetPetaById(w http.ResponseWriter, r *http.Request) {
	idPeta := r.URL.Query().Get("idPeta")
	db := config.GetDB()
	var Map models.Peta

	err := db.Where("idPeta = ?", idPeta).First(&Map)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(Map)

}

func EditPeta(w http.ResponseWriter, r *http.Request) {
	idPeta := r.URL.Query().Get("idPeta")
	db := config.GetDB()

	var peta models.Peta
	if err := db.Where("id = ?", idPeta).First(&peta).Error; err != nil {
		http.Error(w, "Data tidak ditemukan: "+err.Error(), http.StatusNotFound)
		return
	}

	// Decode data baru dari body JSON
	var input models.Peta
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Format data salah: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Update field
	peta.Blok = input.Blok
	peta.Afdeling = input.Afdeling
	peta.Lokasi = input.Lokasi
	peta.Luas = input.Luas
	peta.TahunTanam = input.TahunTanam

	if err := db.Save(&peta).Error; err != nil {
		http.Error(w, "Gagal update data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(peta)
}
