package controllers

import (
	"app-inputan-ptpn/config"
	"app-inputan-ptpn/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// GET ALL
func GetAllPenyadap(w http.ResponseWriter, r *http.Request) {
	var penyadaps []models.Penyadap

	if err := config.DB.Order("created_at desc").Find(&penyadaps).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal mengambil data penyadap: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Data berhasil diambil",
		Data:    penyadaps,
	})
}

// CREATE
func CreatePenyadap(w http.ResponseWriter, r *http.Request) {
	var penyadap models.Penyadap

	if err := json.NewDecoder(r.Body).Decode(&penyadap); err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Format JSON tidak valid: " + err.Error(),
		})
		return
	}

	// Validasi sederhana
	if penyadap.NIK == "" || penyadap.NamaPenyadap == "" {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "ID Mandor, NIK, dan Nama Penyadap wajib diisi",
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

// GET BY NAME (query param ?nama=)
func GetPenyadapByName(w http.ResponseWriter, r *http.Request) {
	nama := r.URL.Query().Get("nama")
	if nama == "" {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Parameter nama wajib diisi",
		})
		return
	}

	var penyadaps []models.Penyadap
	if err := config.DB.Where("nama_penyadap LIKE ?", "%"+nama+"%").Find(&penyadaps).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal mencari penyadap: " + err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Data penyadap ditemukan",
		Data:    penyadaps,
	})
}

// UPDATE
func UpdatePenyadap(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

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

	if err := config.DB.Model(&models.Penyadap{}).Where("id = ?", id).Updates(updates).Error; err != nil {
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

// DELETE
func DeletePenyadap(w http.ResponseWriter, r *http.Request) {
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
