package controllers

import (
	"app-inputan-ptpn/config"
	"app-inputan-ptpn/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

// Struct untuk menerima JSON yang fleksibel
type PetaInput struct {
	Blok        string      `json:"Blok"`
	Code        string      `json:"Code"`
	Afdeling    string      `json:"Afdeling"`
	Luas        float32     `json:"Luas"`
	JumlahPohon int64       `json:"JumlahPohon"`
	JenisKebun  string      `json:"JenisKebun"`
	TahunTanam  interface{} `json:"TahunTanam"` // Terima number atau string
	Kloon       string      `json:"Kloon"`
}

// Helper untuk convert TahunTanam ke string
func convertTahunTanam(value interface{}) string {
	if value == nil {
		return ""
	}

	switch v := value.(type) {
	case string:
		return v
	case float64:
		if v == 0 {
			return ""
		}
		return strconv.Itoa(int(v))
	case int:
		if v == 0 {
			return ""
		}
		return strconv.Itoa(v)
	default:
		return ""
	}
}

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
		log.Printf("Error GetPetaByCode: %v", err)
		http.Error(w, "Data tidak ditemukan untuk code "+code, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(peta)
}

func EditPeta(w http.ResponseWriter, r *http.Request) {
	log.Printf("========== EditPeta Called ==========")
	log.Printf("Method: %s", r.Method)
	log.Printf("URL Path: %s", r.URL.Path)

	vars := mux.Vars(r)
	log.Printf("mux.Vars: %+v", vars)

	idPeta, exists := vars["id"]
	if !exists || idPeta == "" {
		log.Printf("ERROR: Parameter 'id' tidak ditemukan di path")
		http.Error(w, "Parameter 'id' tidak ditemukan", http.StatusBadRequest)
		return
	}

	log.Printf("ID Peta: %s", idPeta)

	db := config.GetDB()

	var peta models.Peta
	if err := db.Where("id = ?", idPeta).First(&peta).Error; err != nil {
		log.Printf("ERROR: Data tidak ditemukan untuk id %s: %v", idPeta, err)
		http.Error(w, "Data tidak ditemukan: "+err.Error(), http.StatusNotFound)
		return
	}

	log.Printf("Data ditemukan: %+v", peta)

	// UBAH: Gunakan PetaInput yang fleksibel
	var input PetaInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Printf("ERROR: Gagal decode JSON: %v", err)
		http.Error(w, "Format data salah: "+err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Input data: %+v", input)

	// Validasi field wajib
	if input.Code == "" || input.Afdeling == "" {
		log.Printf("ERROR: Validasi gagal - Code atau Afdeling kosong")
		http.Error(w, "Field 'Code' dan 'Afdeling' wajib diisi", http.StatusBadRequest)
		return
	}

	// Update field dengan konversi TahunTanam
	peta.Blok = input.Blok
	peta.Code = input.Code
	peta.Afdeling = input.Afdeling
	peta.Luas = input.Luas
	peta.JumlahPohon = input.JumlahPohon
	peta.JenisKebun = input.JenisKebun
	peta.TahunTanam = convertTahunTanam(input.TahunTanam) // PENTING: Convert di sini
	peta.Kloon = input.Kloon

	log.Printf("Data setelah update: %+v", peta)

	if err := db.Save(&peta).Error; err != nil {
		log.Printf("ERROR: Gagal save ke database: %v", err)
		http.Error(w, "Gagal update data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("SUCCESS: Data berhasil disimpan untuk id %s", idPeta)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Data berhasil diperbarui",
		"data":    peta,
	})
}

func UpdatePetaByCode(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Parameter 'code' wajib diisi", http.StatusBadRequest)
		return
	}

	// UBAH: Gunakan PetaInput
	var input PetaInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "Format JSON tidak valid: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validasi field wajib
	if input.Afdeling == "" {
		http.Error(w, "Field 'Afdeling' wajib diisi", http.StatusBadRequest)
		return
	}

	db := config.GetDB()
	var existing models.Peta
	if err := db.Where("code = ?", code).First(&existing).Error; err != nil {
		log.Printf("Error finding peta by code %s: %v", code, err)
		http.Error(w, "Data tidak ditemukan untuk code "+code, http.StatusNotFound)
		return
	}

	// Update semua field dengan konversi
	existing.Blok = input.Blok
	existing.Afdeling = input.Afdeling
	existing.Luas = input.Luas
	existing.JumlahPohon = input.JumlahPohon
	existing.JenisKebun = input.JenisKebun
	existing.TahunTanam = convertTahunTanam(input.TahunTanam) // PENTING: Convert di sini
	existing.Kloon = input.Kloon

	if err := db.Save(&existing).Error; err != nil {
		log.Printf("Error saving peta code %s: %v", code, err)
		http.Error(w, "Gagal menyimpan perubahan: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully updated peta with code: %s", code)
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
		log.Printf("Error getting all peta: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(petas)
}

func CreatePeta(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type harus application/json", http.StatusUnsupportedMediaType)
		return
	}

	// UBAH: Gunakan PetaInput
	var input PetaInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "Format JSON tidak valid: "+err.Error(), http.StatusBadRequest)
		return
	}

	if input.Code == "" || input.Afdeling == "" {
		http.Error(w, "Field 'Code' dan 'Afdeling' wajib diisi", http.StatusBadRequest)
		return
	}

	db := config.GetDB()

	// Convert ke models.Peta dengan konversi TahunTanam
	peta := models.Peta{
		Blok:        input.Blok,
		Code:        input.Code,
		Afdeling:    input.Afdeling,
		Luas:        input.Luas,
		JumlahPohon: input.JumlahPohon,
		JenisKebun:  input.JenisKebun,
		TahunTanam:  convertTahunTanam(input.TahunTanam), // PENTING: Convert di sini
		Kloon:       input.Kloon,
	}

	if err := db.Create(&peta).Error; err != nil {
		log.Printf("Error creating peta: %v", err)
		http.Error(w, "Gagal menyimpan data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully created peta with code: %s", input.Code)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Data peta berhasil dibuat",
		"data":    peta,
	})
}
