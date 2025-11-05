package controllers

import (
	"app-inputan-ptpn/config"
	"app-inputan-ptpn/models"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type SearchMandorResponse struct {
	Success    bool           `json:"success"`
	Message    string         `json:"message"`
	NamaMandor string         `json:"nama_mandor,omitempty"`
	NIK        string         `json:"nik,omitempty"`
	TahunTanam string         `json:"tahun_tanam,omitempty"`
	Data       []models.Rekap `json:"data"`
	Summary    *SummaryMandor `json:"summary,omitempty"`
}

type SummaryMandor struct {
	TotalRecords             int     `json:"total_records"`
	TotalHKO                 int     `json:"total_hko"`
	TotalBasahLatekKebun     float64 `json:"total_basah_latek_kebun"`
	TotalBasahLatekPabrik    float64 `json:"total_basah_latek_pabrik"`
	TotalPersenLatek         float64 `json:"total_persen_latek"`
	TotalBasahLumpKebun      float64 `json:"total_basah_lump_kebun"`
	TotalBasahLumpPabrik     float64 `json:"total_basah_lump_pabrik"`
	TotalPersenLump          float64 `json:"total_persen_lump"`
	TotalK3Sheet             float64 `json:"total_k3_sheet"`
	TotalKeringSheet         float64 `json:"total_kering_sheet"`
	TotalKeringBrCr          float64 `json:"total_kering_br_cr"`
	TotalKeringJumlah        float64 `json:"total_kering_jumlah"`
	RataRataProduksiPerTaper float64 `json:"rata_rata_produksi_per_taper"`
}

// SearchMandor mencari data rekap produksi berdasarkan parameter mandor
// Parameter:
// - idMandor (wajib): ID mandor
// - tanggalAwal (opsional): tanggal mulai pencarian
// - tanggalAkhir (opsional): tanggal akhir pencarian
// - tipeProduksi (opsional): jenis tipe produksi
func SearchMandor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idMandorStr := r.URL.Query().Get("idMandor")
	tanggalAwal := r.URL.Query().Get("tanggalAwal")
	tanggalAkhir := r.URL.Query().Get("tanggalAkhir")
	tipeProduksi := r.URL.Query().Get("tipeProduksi")

	if idMandorStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(SearchMandorResponse{
			Success: false,
			Message: "Parameter idMandor diperlukan",
		})
		return
	}

	// Convert idMandor string to int
	idMandor, err := strconv.Atoi(idMandorStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(SearchMandorResponse{
			Success: false,
			Message: "Format idMandor tidak valid, harus berupa angka",
		})
		return
	}

	// Jika ada tanggal akhir tapi tidak ada tanggal awal
	if tanggalAwal == "" && tanggalAkhir != "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(SearchMandorResponse{
			Success: false,
			Message: "Tanggal awal diperlukan jika menggunakan tanggal akhir",
		})
		return
	}

	// STEP 1: Get mandor data menggunakan GetMandor
	namaMandor, nik, tahunTanam, err := GetMandor(idMandor)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(SearchMandorResponse{
			Success: false,
			Message: "Mandor dengan ID tersebut tidak ditemukan",
		})
		return
	}

	// STEP 2: Cari data rekap berdasarkan NIK yang didapat
	db := config.GetDB()
	var rekapList []models.Rekap
	query := db.Model(&models.Rekap{}).Where("nik = ?", nik).Where("tipe_produksi != ?", "REKAPITULASI")

	// Jika hanya ada idMandor
	if tanggalAwal == "" && tanggalAkhir == "" && tipeProduksi == "" {
		if err := query.Order("tanggal DESC").Find(&rekapList).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(SearchMandorResponse{
				Success: false,
				Message: "Gagal mengambil data: " + err.Error(),
			})
			return
		}
	}

	// Jika hanya ada tanggal awal saja
	if tanggalAwal != "" && tanggalAkhir == "" && tipeProduksi == "" {
		tglAwal, err := time.Parse("2006-01-02", tanggalAwal)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(SearchMandorResponse{
				Success: false,
				Message: "Format tanggal awal tidak valid (gunakan: YYYY-MM-DD)",
			})
			return
		}

		query = query.Where("tanggal >= ?", tglAwal)
		if err := query.Order("tanggal DESC").Find(&rekapList).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(SearchMandorResponse{
				Success: false,
				Message: "Gagal mengambil data: " + err.Error(),
			})
			return
		}
	}

	// Jika hanya ada tanggal awal dan akhir saja
	if tanggalAwal != "" && tanggalAkhir != "" && tipeProduksi == "" {
		tglAwal, err := time.Parse("2006-01-02", tanggalAwal)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(SearchMandorResponse{
				Success: false,
				Message: "Format tanggal awal tidak valid (gunakan: YYYY-MM-DD)",
			})
			return
		}

		tglAkhir, err := time.Parse("2006-01-02", tanggalAkhir)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(SearchMandorResponse{
				Success: false,
				Message: "Format tanggal akhir tidak valid (gunakan: YYYY-MM-DD)",
			})
			return
		}

		query = query.Where("tanggal BETWEEN ? AND ?", tglAwal, tglAkhir)
		if err := query.Order("tanggal DESC").Find(&rekapList).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(SearchMandorResponse{
				Success: false,
				Message: "Gagal mengambil data: " + err.Error(),
			})
			return
		}
	}

	// Jika hanya ada tipe produksi
	if tanggalAwal == "" && tanggalAkhir == "" && tipeProduksi != "" {
		query = query.Where("tipe_produksi = ?", tipeProduksi)
		if err := query.Order("tanggal DESC").Find(&rekapList).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(SearchMandorResponse{
				Success: false,
				Message: "Gagal mengambil data: " + err.Error(),
			})
			return
		}
	}

	// Jika ada semua parameter
	if tipeProduksi != "" && tanggalAwal != "" && tanggalAkhir != "" {
		tglAwal, err := time.Parse("2006-01-02", tanggalAwal)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(SearchMandorResponse{
				Success: false,
				Message: "Format tanggal awal tidak valid (gunakan: YYYY-MM-DD)",
			})
			return
		}

		tglAkhir, err := time.Parse("2006-01-02", tanggalAkhir)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(SearchMandorResponse{
				Success: false,
				Message: "Format tanggal akhir tidak valid (gunakan: YYYY-MM-DD)",
			})
			return
		}

		query = query.Where("tanggal BETWEEN ? AND ?", tglAwal, tglAkhir).
			Where("tipe_produksi = ?", tipeProduksi)

		if err := query.Order("tanggal DESC").Find(&rekapList).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(SearchMandorResponse{
				Success: false,
				Message: "Gagal mengambil data: " + err.Error(),
			})
			return
		}
	}

	// Hitung summary
	summary := calculateMandorSummary(rekapList)

	// Response success dengan informasi mandor
	json.NewEncoder(w).Encode(SearchMandorResponse{
		Success:    true,
		Message:    "Data berhasil diambil",
		NamaMandor: namaMandor,
		NIK:        nik,
		TahunTanam: tahunTanam,
		Data:       rekapList,
		Summary:    summary,
	})
}

// calculateMandorSummary menghitung total dari data rekap mandor
func calculateMandorSummary(data []models.Rekap) *SummaryMandor {
	if len(data) == 0 {
		return nil
	}

	summary := &SummaryMandor{
		TotalRecords: len(data),
	}

	var totalProduksiPerTaper float64
	var countProduksiPerTaper int

	for _, r := range data {
		summary.TotalHKO += r.HKOHariIni
		summary.TotalBasahLatekKebun += r.HariIniBasahLatekKebun
		summary.TotalBasahLatekPabrik += r.HariIniBasahLatekPabrik
		summary.TotalBasahLumpKebun += r.HariIniBasahLumpKebun
		summary.TotalBasahLumpPabrik += r.HariIniBasahLumpPabrik
		summary.TotalKeringSheet += r.HariIniKeringSheet
		summary.TotalKeringBrCr += r.HariIniKeringBrCr
		summary.TotalKeringJumlah += r.HariIniKeringJumlah

		// Hitung rata-rata produksi per taper
		if r.ProduksiPerTaperHariIni > 0 {
			totalProduksiPerTaper += r.ProduksiPerTaperHariIni
			countProduksiPerTaper++
		}
	}
	//hitung k3
	summary.TotalK3Sheet = summary.TotalKeringSheet / summary.TotalBasahLatekPabrik
	//hitung persen latek
	summary.TotalPersenLatek = (summary.TotalBasahLatekKebun - summary.TotalBasahLatekPabrik) / summary.TotalBasahLumpKebun * 100
	//hitung persen lump
	summary.TotalPersenLump = (summary.TotalBasahLumpKebun - summary.TotalBasahLumpPabrik) / summary.TotalBasahLumpPabrik * 100
	// Hitung rata-rata produksi per taper
	if countProduksiPerTaper > 0 {
		summary.RataRataProduksiPerTaper = totalProduksiPerTaper / float64(countProduksiPerTaper)
	}

	return summary
}

// GetMandor mengambil data mandor berdasarkan ID
// Returns: namaMandor, NIK, tahunTanam, error
func GetMandor(idMandor int) (string, string, string, error) {
	db := config.GetDB()
	mandor := models.Mandor{}

	if err := db.First(&mandor, idMandor).Error; err != nil {
		return "", "", "", err
	}

	return mandor.Nama, mandor.NIK, mandor.TahunTanam, nil
}
