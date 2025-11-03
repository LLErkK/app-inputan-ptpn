package controllers

import (
	"app-inputan-ptpn/config"
	"app-inputan-ptpn/models"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

// CreateMaster membuat record Master baru dan mengembalikan ID-nya
func CreateMaster(tanggal time.Time, afdeling string, namaFile string) (uint64, error) {
	db := config.GetDB()

	master := models.Master{
		Tanggal:  tanggal,
		Afdeling: afdeling,
		NamaFile: namaFile,
	}

	// Simpan ke database
	if err := db.Create(&master).Error; err != nil {
		return 0, err
	}

	// Kembalikan ID master yang baru dibuat
	return master.ID, nil
}

// GetAllMaster mengembalikan semua master dalam bentuk JSON
func GetAllMaster(w http.ResponseWriter, r *http.Request) {
	db := config.GetDB()
	var masters []models.Master

	// Load semua master beserta relasinya
	if err := db.Preload("Rekaps").Preload("Produksis").Find(&masters).Error; err != nil {
		http.Error(w, "Gagal mengambil data master: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(masters); err != nil {
		http.Error(w, "Gagal encode JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// DeleteMaster menghapus master beserta semua Rekap dan Produksi terkait (cascade)
func DeleteMaster(w http.ResponseWriter, r *http.Request) {
	db := config.GetDB()
	vars := mux.Vars(r)
	idStr, ok := vars["masterId"]
	if !ok {
		http.Error(w, "ID master tidak ditemukan di URL", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "ID master tidak valid", http.StatusBadRequest)
		return
	}

	var master models.Master
	if err := db.First(&master, id).Error; err != nil {
		http.Error(w, fmt.Sprintf("Master dengan ID %d tidak ditemukan", id), http.StatusNotFound)
		return
	}

	if err := db.Delete(&master).Error; err != nil {
		http.Error(w, fmt.Sprintf("Gagal menghapus master: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Master dengan ID %d berhasil dihapus", id)))
}
