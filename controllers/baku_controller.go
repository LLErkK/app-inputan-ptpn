package controllers

import (
	"app-inputan-ptpn/config"
	"app-inputan-ptpn/models"
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func respondJSON(w http.ResponseWriter, status int, payload APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// ===================== GET ALL =====================
func GetAllBakuPenyadap(w http.ResponseWriter, r *http.Request) {
	var penyadap []models.BakuPenyadap

	if err := config.DB.Order("created_at desc").Find(&penyadap).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal mengambil data: " + err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Data berhasil diambil",
		Data:    penyadap,
	})
}

func GetAllBakuMandor(w http.ResponseWriter, r *http.Request) {
	var mandor []models.BakuMandor
	if err := config.DB.Order("created_at desc").Find(&mandor).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal mengambil data: " + err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Data berhasil diambil",
		Data:    mandor,
	})
}

// ===================== GET BY ID =====================
func GetBakuPenyadapByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var penyadap models.BakuPenyadap
	if err := config.DB.First(&penyadap, id).Error; err != nil {
		respondJSON(w, http.StatusNotFound, APIResponse{
			Success: false,
			Message: "Data tidak ditemukan",
		})
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Data berhasil ditemukan",
		Data:    penyadap,
	})
}

// ===================== CREATE =====================
func CreateBakuPenyadap(w http.ResponseWriter, r *http.Request) {
	var penyadap models.BakuPenyadap
	if err := json.NewDecoder(r.Body).Decode(&penyadap); err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Format JSON tidak valid: " + err.Error(),
		})
		return
	}

	// Validasi sederhana
	if penyadap.IdBakuMandor == 0 || penyadap.NIK == "" || penyadap.NamaPenyadap == "" {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "ID mandor, NIK, dan Nama penyadap wajib diisi",
		})
		return
	}

	if err := config.DB.Create(&penyadap).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal menyimpan data penyadap: " + err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusCreated, APIResponse{
		Success: true,
		Message: "Data penyadap berhasil ditambahkan",
		Data:    penyadap,
	})
}

// ===================== UPDATE =====================
func UpdateBakuPenyadap(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	// pastikan id valid
	_, err := strconv.Atoi(id)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "ID tidak valid",
		})
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Format JSON tidak valid: " + err.Error(),
		})
		return
	}

	if err := config.DB.Model(&models.BakuPenyadap{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal mengupdate data penyadap: " + err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Data penyadap berhasil diperbarui",
	})
}

// ===================== DELETE =====================
func DeleteBakuPenyadap(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if err := config.DB.Delete(&models.BakuPenyadap{}, id).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal menghapus data penyadap: " + err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Data penyadap berhasil dihapus",
	})
}

type BakuPageData struct {
	Title        string
	MandorList   []models.BakuMandor
	PenyadapList []models.BakuPenyadap
}

// Tambahkan ini di bagian atas file baku_controller.go setelah import

// Template functions
var templateFuncs = template.FuncMap{
	"add": func(a, b int) int {
		return a + b
	},
}

// Update fungsi ServeBakuPage
func ServeBakuPage(w http.ResponseWriter, r *http.Request) {
	var mandor []models.BakuMandor
	var penyadap []models.BakuPenyadap

	// Ambil data dengan relasi
	if err := config.DB.Order("created_at desc").Find(&mandor).Error; err != nil {
		http.Error(w, "Gagal mengambil data mandor: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Preload relasi Mandor untuk penyadap
	if err := config.DB.Preload("Mandor").Order("created_at desc").Find(&penyadap).Error; err != nil {
		http.Error(w, "Gagal mengambil data penyadap: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Data yang dikirim ke template
	data := BakuPageData{
		Title:        "Data Mandor & Penyadap",
		MandorList:   mandor,
		PenyadapList: penyadap,
	}

	// Parse file template dengan custom functions
	tmpl, err := template.New("baku.html").Funcs(templateFuncs).ParseFiles("templates/html/baku.html")
	if err != nil {
		http.Error(w, "Gagal parse template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Render template dengan data
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Gagal render template: "+err.Error(), http.StatusInternalServerError)
	}
}
