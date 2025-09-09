package controllers

import (
	"app-inputan-ptpn/config"
	"app-inputan-ptpn/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func GetAllMandor(w http.ResponseWriter, r *http.Request) {
	var mandors []models.BakuMandor

	if err := config.DB.Order("created_at desc").Find(&mandors).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal mengambil data mandor: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Data berhasil diambil",
		Data:    mandors,
	})
}

func CreateMandor(w http.ResponseWriter, r *http.Request) {
	var mandor models.BakuMandor
	if err := json.NewDecoder(r.Body).Decode(&mandor); err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Format JSON tidak valid: " + err.Error(),
		})
		return
	}
	if err := config.DB.Create(&mandor).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal menyimpan data mandor: " + err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusCreated, APIResponse{
		Success: true,
		Message: "Data mandor berhasil ditambahkan",
		Data:    mandor,
	})

}

func UpdateMandor(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	_, err := strconv.Atoi(id)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "ID tidak valid",
		})
		return
	}

	var update models.BakuMandor
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Format JSON tidak valid: " + err.Error(),
		})
		return
	}

	// update hanya field yang dikirim
	if err := config.DB.Model(&models.BakuMandor{}).
		Where("id = ?", id).
		Updates(update).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal menyimpan data mandor: " + err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Data mandor berhasil diperbarui",
	})
}

func DeleteMandor(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	_, err := strconv.Atoi(id)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "ID tidak valid",
		})
		return
	}

	if err := config.DB.Delete(&models.BakuMandor{}, id).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal menghapus data mandor: " + err.Error(),
		})
	}
	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Data mandor berhasil ditambahkan",
	})
}
