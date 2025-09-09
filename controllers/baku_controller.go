package controllers

import (
	"app-inputan-ptpn/config"
	"app-inputan-ptpn/models"
	"encoding/json"
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Response standar
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

	if err := config.DB.Preload("Mandor").Preload("Penyadap").Order("created_at desc").Find(&penyadap).Error; err != nil {
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
	if err := config.DB.Preload("Mandor").Preload("Penyadap").First(&penyadap, id).Error; err != nil {
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

	// Validasi
	if penyadap.IdBakuMandor == 0 || penyadap.IdPenyadap == 0 {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "ID mandor dan ID penyadap wajib diisi",
		})
		return
	}
	if penyadap.Tanggal.IsZero() {
		penyadap.Tanggal = time.Now()
	}

	// Simpan data penyadap
	if err := config.DB.Create(&penyadap).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal menyimpan data penyadap: " + err.Error(),
		})
		return
	}

	// Update detail harian
	updateBakuDetail(penyadap, "create", nil)

	respondJSON(w, http.StatusCreated, APIResponse{
		Success: true,
		Message: "Data penyadap berhasil ditambahkan",
		Data:    penyadap,
	})
}

// ===================== UPDATE =====================
func UpdateBakuPenyadap(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var existing models.BakuPenyadap
	if err := config.DB.First(&existing, id).Error; err != nil {
		respondJSON(w, http.StatusNotFound, APIResponse{
			Success: false,
			Message: "Data penyadap tidak ditemukan",
		})
		return
	}

	var updates models.BakuPenyadap
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Format JSON tidak valid: " + err.Error(),
		})
		return
	}

	// Simpan selisih untuk update detail
	oldCopy := existing

	if err := config.DB.Model(&existing).Updates(updates).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal update penyadap: " + err.Error(),
		})
		return
	}

	// Update detail
	updateBakuDetail(existing, "update", &oldCopy)

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Data penyadap berhasil diperbarui",
	})
}

// ===================== DELETE =====================
func DeleteBakuPenyadap(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var penyadap models.BakuPenyadap
	if err := config.DB.First(&penyadap, id).Error; err != nil {
		respondJSON(w, http.StatusNotFound, APIResponse{
			Success: false,
			Message: "Data penyadap tidak ditemukan",
		})
		return
	}

	if err := config.DB.Delete(&penyadap).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal menghapus data penyadap: " + err.Error(),
		})
		return
	}

	updateBakuDetail(penyadap, "delete", nil)

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Data penyadap berhasil dihapus",
	})
}

// ===================== DETAIL UPDATER =====================
func updateBakuDetail(entry models.BakuPenyadap, action string, oldEntry *models.BakuPenyadap) {
	var detail models.BakuDetail
	err := config.DB.Where("tanggal = ?", entry.Tanggal).First(&detail).Error
	if err != nil {
		// kalau belum ada & action create → buat baru
		if action == "create" {
			detail = models.BakuDetail{Tanggal: entry.Tanggal}
		} else {
			return
		}
	}

	switch action {
	case "create":
		detail.JumlahKebunBasahLatek += entry.BasahLatex
		detail.JumlahSheet += entry.Sheet
		detail.JumlahKebunBasahLump += entry.BasahLump
		detail.JumlahBrCr += entry.BrCr
	case "update":
		if oldEntry != nil {
			detail.JumlahKebunBasahLatek += entry.BasahLatex - oldEntry.BasahLatex
			detail.JumlahSheet += entry.Sheet - oldEntry.Sheet
			detail.JumlahKebunBasahLump += entry.BasahLump - oldEntry.BasahLump
			detail.JumlahBrCr += entry.BrCr - oldEntry.BrCr
		}
	case "delete":
		detail.JumlahKebunBasahLatek -= entry.BasahLatex
		detail.JumlahSheet -= entry.Sheet
		detail.JumlahKebunBasahLump -= entry.BasahLump
		detail.JumlahBrCr -= entry.BrCr
	}

	config.DB.Save(&detail)
}

// ===================== PAGE =====================
type BakuPageData struct {
	Title        string
	MandorList   []models.BakuMandor
	PenyadapList []models.BakuPenyadap
}

// Template functions
var templateFuncs = template.FuncMap{
	"add": func(a, b int) int { return a + b },
}

func ServeBakuPage(w http.ResponseWriter, r *http.Request) {
	var mandor []models.BakuMandor
	var penyadap []models.BakuPenyadap

	if err := config.DB.Order("created_at desc").Find(&mandor).Error; err != nil {
		http.Error(w, "Gagal mengambil data mandor: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := config.DB.Preload("Mandor").Preload("Penyadap").Order("created_at desc").Find(&penyadap).Error; err != nil {
		http.Error(w, "Gagal mengambil data penyadap: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := BakuPageData{
		Title:        "Data Mandor & Penyadap",
		MandorList:   mandor,
		PenyadapList: penyadap,
	}

	tmpl, err := template.New("baku.html").Funcs(templateFuncs).ParseFiles("templates/html/baku.html")
	if err != nil {
		http.Error(w, "Gagal parse template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Gagal render template: "+err.Error(), http.StatusInternalServerError)
	}
}

// ===================== DETAIL =====================
func GetAllBakuDetail(w http.ResponseWriter, r *http.Request) {
	var details []models.BakuDetail
	if err := config.DB.Order("tanggal desc").Find(&details).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal mengambil detail baku: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Data detail berhasil diambil",
		Data:    details,
	})
}

func GetBakuDetailByDate(w http.ResponseWriter, r *http.Request) {
	tanggal := mux.Vars(r)["tanggal"]

	var detail models.BakuDetail
	if err := config.DB.Where("tanggal = ?", tanggal).First(&detail).Error; err != nil {
		respondJSON(w, http.StatusNotFound, APIResponse{
			Success: false,
			Message: "Detail untuk tanggal " + tanggal + " tidak ditemukan",
		})
		return
	}
	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Detail berhasil ditemukan",
		Data:    detail,
	})
}
