package controllers

import (
	"app-inputan-ptpn/config"
	"app-inputan-ptpn/models"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type SearchResponse struct {
	Success      bool              `json:"success"`
	Message      string            `json:"message"`
	NamaPenyadap string            `json:"nama_penyadap,omitempty"`
	NIK          string            `json:"nik,omitempty"`
	Data         []models.Produksi `json:"data"`
	Summary      *SummaryData      `json:"summary,omitempty"`
}

type SummaryData struct {
	TotalRecords int     `json:"total_records"`
	TotalLatek   float64 `json:"total_basah_latek"`
	TotalSheet   float64 `json:"total_sheet"`
	TotalLump    float64 `json:"total_basah_lump"`
	TotalBrCr    float64 `json:"total_br_cr"`
}

// SearchPenyadap mencari data rekap produksi berdasarkan parameter
// Parameter:
// - idPenyadap (wajib): ID penyadap
// - tanggalAwal (opsional): tanggal mulai pencarian
// - tanggalAkhir (opsional): tanggal akhir pencarian
// - tipeProduksi (opsional): jenis tipe produksi
// - afdeling (opsional): filter berdasarkan afdeling
func SearchPenyadap(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idPenyadapStr := r.URL.Query().Get("idPenyadap")
	tanggalAwal := r.URL.Query().Get("tanggalAwal")
	tanggalAkhir := r.URL.Query().Get("tanggalAkhir")
	tipeProduksi := r.URL.Query().Get("tipeProduksi")
	afdeling := r.URL.Query().Get("afdeling")

	if idPenyadapStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(SearchResponse{
			Success: false,
			Message: "Parameter idPenyadap diperlukan",
		})
		return
	}

	// Convert idPenyadap string to int
	idPenyadap, err := strconv.Atoi(idPenyadapStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(SearchResponse{
			Success: false,
			Message: "Format idPenyadap tidak valid, harus berupa angka",
		})
		return
	}

	// Jika ada tanggal akhir tapi tidak ada tanggal awal
	if tanggalAwal == "" && tanggalAkhir != "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(SearchResponse{
			Success: false,
			Message: "Tanggal awal diperlukan jika menggunakan tanggal akhir",
		})
		return
	}

	// STEP 1: Get penyadap data menggunakan GetPenyadap
	namaPenyadap, nik, err := GetPenyadap(idPenyadap)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(SearchResponse{
			Success: false,
			Message: "Penyadap dengan ID tersebut tidak ditemukan",
		})
		return
	}

	// STEP 2: Cari data produksi berdasarkan NIK yang didapat
	db := config.GetDB()
	var produksiList []models.Produksi
	// Base query dengan filter NIK
	query := db.Model(&models.Produksi{}).Where("nik = ?", nik)

	// Filter berdasarkan afdeling jika ada
	if afdeling != "" {
		query = query.Where("afdeling = ?", afdeling)
	}

	// Filter berdasarkan tanggal
	if tanggalAwal != "" {
		tglAwal, err := time.Parse("2006-01-02", tanggalAwal)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(SearchResponse{
				Success: false,
				Message: "Format tanggal awal tidak valid (gunakan: YYYY-MM-DD)",
			})
			return
		}
		tglAwal = tglAwal.AddDate(0, 0, -1)

		if tanggalAkhir != "" {
			tglAkhir, err := time.Parse("2006-01-02", tanggalAkhir)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(SearchResponse{
					Success: false,
					Message: "Format tanggal akhir tidak valid (gunakan: YYYY-MM-DD)",
				})
				return
			}

			query = query.Where("tanggal BETWEEN ? AND ?", tglAwal, tglAkhir)
		} else {
			query = query.Where("tanggal >= ?", tglAwal)
		}
	}

	// Filter berdasarkan tipe produksi jika ada
	if tipeProduksi != "" {
		query = query.Where("tipe_produksi = ?", tipeProduksi)
	}

	// Execute query
	if err := query.Order("tanggal DESC").Find(&produksiList).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(SearchResponse{
			Success: false,
			Message: "Gagal mengambil data: " + err.Error(),
		})
		return
	}

	// Hitung summary
	summary := calculateSummary(produksiList)

	// Response success dengan informasi penyadap
	json.NewEncoder(w).Encode(SearchResponse{
		Success:      true,
		Message:      "Data berhasil diambil",
		NamaPenyadap: namaPenyadap,
		NIK:          nik,
		Data:         produksiList,
		Summary:      summary,
	})
}

// calculateSummary menghitung total dari data produksi
func calculateSummary(data []models.Produksi) *SummaryData {
	if len(data) == 0 {
		return nil
	}

	summary := &SummaryData{
		TotalRecords: len(data),
	}

	for _, p := range data {
		summary.TotalLatek += p.BasahLatek
		summary.TotalSheet += p.Sheet
		summary.TotalLump += p.BasahLump
		summary.TotalBrCr += p.BrCr
	}

	return summary
}

// GetPenyadap mengambil data penyadap berdasarkan ID
func GetPenyadap(idPenyadap int) (string, string, error) {
	db := config.GetDB()
	penyadap := models.Penyadap{}

	if err := db.First(&penyadap, idPenyadap).Error; err != nil {
		return "", "", err
	}

	return penyadap.NamaPenyadap, penyadap.NIK, nil
}
