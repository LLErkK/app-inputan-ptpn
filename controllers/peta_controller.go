package controllers

import (
	"app-inputan-ptpn/config"
	"app-inputan-ptpn/models"
	"encoding/json"
	"net/http"
)

func ServePetaPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/html/peta.html")
}
func GetPetaByCode(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Parameter 'code' wajib diisi", http.StatusBadRequest)
		return
	}

	db := config.GetDB()
	var peta models.Peta
	if err := db.Where("code = ?", code).First(&peta).Error; err != nil {
		http.Error(w, "Data tidak ditemukan untuk code "+code, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(peta)
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

	// Validasi field wajib
	if input.Code == "" || input.Afdeling == "" {
		http.Error(w, "Field 'Code' dan 'Afdeling' wajib diisi", http.StatusBadRequest)
		return
	}

	// Update field
	peta.Blok = input.Blok
	peta.Code = input.Code
	peta.Afdeling = input.Afdeling
	peta.Luas = input.Luas
	peta.JumlahPohon = input.JumlahPohon
	peta.JenisKebun = input.JenisKebun
	peta.TahunTanam = input.TahunTanam
	peta.Kloon = input.Kloon

	if err := db.Save(&peta).Error; err != nil {
		http.Error(w, "Gagal update data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(peta)
}

func UpdatePetaByCode(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Parameter 'code' wajib diisi", http.StatusBadRequest)
		return
	}

	var updatedData models.Peta
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		http.Error(w, "Format JSON tidak valid: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validasi field wajib
	if updatedData.Afdeling == "" {
		http.Error(w, "Field 'Afdeling' wajib diisi", http.StatusBadRequest)
		return
	}

	db := config.GetDB()
	var existing models.Peta
	if err := db.Where("code = ?", code).First(&existing).Error; err != nil {
		http.Error(w, "Data tidak ditemukan untuk code "+code, http.StatusNotFound)
		return
	}

	// Update field (code tidak diupdate karena digunakan sebagai identifier)
	existing.Blok = updatedData.Blok
	existing.Afdeling = updatedData.Afdeling
	existing.Luas = updatedData.Luas
	existing.JumlahPohon = updatedData.JumlahPohon
	existing.JenisKebun = updatedData.JenisKebun
	existing.TahunTanam = updatedData.TahunTanam

	if err := db.Save(&existing).Error; err != nil {
		http.Error(w, "Gagal menyimpan perubahan: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Data berhasil diperbarui",
		"data":    existing,
	})
}

func GetAllPeta(w http.ResponseWriter, r *http.Request) {
	db := config.GetDB()
	var petas []models.Peta
	if err := db.Find(&petas).Error; err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(petas)
}

func CreatePeta(w http.ResponseWriter, r *http.Request) {
	// Pastikan request body berupa JSON
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type harus application/json", http.StatusUnsupportedMediaType)
		return
	}

	// Decode body ke struct models.Peta
	var input models.Peta
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Format JSON tidak valid: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validasi hanya field yang wajib (Code dan Afdeling)
	if input.Code == "" || input.Afdeling == "" {
		http.Error(w, "Field 'Code' dan 'Afdeling' wajib diisi", http.StatusBadRequest)
		return
	}

	db := config.GetDB()

	// Simpan ke database
	if err := db.Create(&input).Error; err != nil {
		http.Error(w, "Gagal menyimpan data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Beri response sukses
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Data peta berhasil dibuat",
		"data":    input,
	})
}
