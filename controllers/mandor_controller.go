// controllers/mandor_controller.go - Updated with Tipe validation

package controllers

import (
	"app-inputan-ptpn/config"
	"app-inputan-ptpn/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
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

// UPDATED: CreateMandor - Now validates and sets Tipe
func CreateMandor(w http.ResponseWriter, r *http.Request) {
	var mandor models.BakuMandor
	if err := json.NewDecoder(r.Body).Decode(&mandor); err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Format JSON tidak valid: " + err.Error(),
		})
		return
	}

	// UPDATED: Validate required fields including tipe
	if mandor.Mandor == "" || mandor.Afdeling == "" {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Nama mandor dan afdeling wajib diisi",
		})
		return
	}

	// UPDATED: Validate tipe produksi
	if mandor.Tipe != "" && !models.IsValidTipeProduksi(mandor.Tipe) {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Tipe produksi tidak valid. Pilih: BAKU, BAKU_BORONG, BORONG_EXTERNAL, BORONG_INTERNAL, TETES_LANJUT, atau BORONG_MINGGU",
		})
		return
	}

	// Set default tipe if empty
	if mandor.Tipe == "" {
		mandor.Tipe = models.TipeBaku
	}

	// Set default tahun tanam if empty
	if mandor.TahunTanam == 0 {
		mandor.TahunTanam = 2024
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
		Message: "Data mandor berhasil ditambahkan dengan tipe " + string(mandor.Tipe),
		Data:    mandor,
	})
}

// UPDATED: UpdateMandor - Now handles tipe updates and cascades to existing BakuPenyadap
func UpdateMandor(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	mandorID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "ID tidak valid",
		})
		return
	}

	// Get existing mandor
	var existingMandor models.BakuMandor
	if err := config.DB.First(&existingMandor, mandorID).Error; err != nil {
		respondJSON(w, http.StatusNotFound, APIResponse{
			Success: false,
			Message: "Data mandor tidak ditemukan",
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

	// UPDATED: Validate tipe if provided
	if update.Tipe != "" && !models.IsValidTipeProduksi(update.Tipe) {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Tipe produksi tidak valid",
		})
		return
	}

	// Check if tipe is being changed
	tipeChanged := update.Tipe != "" && update.Tipe != existingMandor.Tipe

	// Update mandor
	if err := config.DB.Model(&existingMandor).Updates(update).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Gagal update mandor: " + err.Error(),
		})
		return
	}

	// UPDATED: If tipe changed, update all related BakuPenyadap records
	if tipeChanged {
		err := config.DB.Model(&models.BakuPenyadap{}).
			Where("id_baku_mandor = ?", mandorID).
			Update("tipe", update.Tipe).Error

		if err != nil {
			// Log the error but don't fail the entire operation
			config.DB.Rollback()
			respondJSON(w, http.StatusInternalServerError, APIResponse{
				Success: false,
				Message: "Gagal update tipe pada data penyadap terkait: " + err.Error(),
			})
			return
		}

		// UPDATED: Recalculate all BakuDetail records for this mandor
		// This is more complex as we need to recalculate based on the new tipe grouping
		if err := recalculateAllBakuDetailForMandor(uint(mandorID), existingMandor.Tipe, update.Tipe); err != nil {
			// Log error but continue
			config.DB.Rollback()
			respondJSON(w, http.StatusInternalServerError, APIResponse{
				Success: false,
				Message: "Gagal update detail summary: " + err.Error(),
			})
			return
		}

		respondJSON(w, http.StatusOK, APIResponse{
			Success: true,
			Message: "Data mandor berhasil diperbarui. Tipe produksi diubah dari " + string(existingMandor.Tipe) + " ke " + string(update.Tipe) + " dan semua data penyadap terkait telah diperbarui.",
		})
	} else {
		respondJSON(w, http.StatusOK, APIResponse{
			Success: true,
			Message: "Data mandor berhasil diperbarui",
		})
	}
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
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Data mandor berhasil dihapus",
	})
}

func GetMandorByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var mandor models.BakuMandor
	if err := config.DB.First(&mandor, id).Error; err != nil {
		respondJSON(w, http.StatusNotFound, APIResponse{
			Success: false,
			Message: "Data mandor tidak ditemukan",
		})
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Data mandor berhasil ditemukan",
		Data:    mandor,
	})
}

// UPDATED: Helper function to recalculate BakuDetail when mandor tipe changes
func recalculateAllBakuDetailForMandor(mandorID uint, oldTipe, newTipe models.TipeProduksi) error {
	// Get mandor info
	var mandor models.BakuMandor
	if err := config.DB.First(&mandor, mandorID).Error; err != nil {
		return err
	}

	// Get all unique dates for this mandor's BakuPenyadap records
	var dates []time.Time
	err := config.DB.Model(&models.BakuPenyadap{}).
		Where("id_baku_mandor = ?", mandorID).
		Distinct("DATE(tanggal)").
		Pluck("DATE(tanggal)", &dates).Error

	if err != nil {
		return err
	}

	// For each date, recalculate the BakuDetail
	for _, date := range dates {
		// Remove old BakuDetail record with old tipe
		config.DB.Where("DATE(tanggal) = DATE(?) AND mandor = ? AND tipe = ?",
			date, mandor.Mandor, oldTipe).Delete(&models.BakuDetail{})

		// Recalculate with new tipe
		if err := RecalculateBakuDetail(date, mandorID, newTipe); err != nil {
			return err
		}
	}

	return nil
}
