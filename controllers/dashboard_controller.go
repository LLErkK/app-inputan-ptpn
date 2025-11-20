package controllers

import (
	"app-inputan-ptpn/config"
	"app-inputan-ptpn/models"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"time"
)

func ServeDashboardPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/html/dashboard.html")
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

type dashboardDataResponse struct {
	// HKO
	TotalHKOHariIni       int `json:"totalHKOHariIni"`
	TotalHKOSampaiHariIni int `json:"totalHKOSampaiHariIni"`

	// Hari Ini - Basah Latek
	TotalHariIniBasahLatekKebun  float64 `json:"totalHariIniBasahLatekKebun"`
	TotalHariIniBasahLatekPabrik float64 `json:"totalHariIniBasahLatekPabrik"`
	TotalHariIniBasahLatekPersen float64 `json:"totalHariIniBasahLatekPersen"` // (Kebun - Pabrik) / Pabrik * 100

	// Hari Ini - Basah Lump
	TotalHariIniBasahLumpKebun  float64 `json:"totalHariIniBasahLumpKebun"`
	TotalHariIniBasahLumpPabrik float64 `json:"totalHariIniBasahLumpPabrik"`
	TotalHariIniBasahLumpPersen float64 `json:"totalHariIniBasahLumpPersen"` // (Kebun - Pabrik) / Pabrik * 100

	// Hari Ini - Kering
	TotalHariIniK3Sheet       float64 `json:"totalHariIniK3Sheet"` // SUM dari K3Sheet (bukan perhitungan)
	TotalHariIniKeringSheet   float64 `json:"totalHariIniKeringSheet"`
	TotalHariIniKeringBrCr    float64 `json:"totalHariIniKeringBrCr"`
	TotalHariIniKeringJumlah  float64 `json:"totalHariIniKeringJumlah"`
	TotalHariIniK3SheetPersen float64 `json:"totalHariIniK3SheetPersen"` // KeringSheet / BasahLatekPabrik * 100

	// Sampai Hari Ini - Basah Latek
	TotalSampaiHariIniBasahLatekKebun  float64 `json:"totalSampaiHariIniBasahLatekKebun"`
	TotalSampaiHariIniBasahLatekPabrik float64 `json:"totalSampaiHariIniBasahLatekPabrik"`
	TotalSampaiHariIniBasahLatekPersen float64 `json:"totalSampaiHariIniBasahLatekPersen"` // (Kebun - Pabrik) / Pabrik * 100

	// Sampai Hari Ini - Basah Lump
	TotalSampaiHariIniBasahLumpKebun  float64 `json:"totalSampaiHariIniBasahLumpKebun"`
	TotalSampaiHariIniBasahLumpPabrik float64 `json:"totalSampaiHariIniBasahLumpPabrik"`
	TotalSampaiHariIniBasahLumpPersen float64 `json:"totalSampaiHariIniBasahLumpPersen"` // (Kebun - Pabrik) / Pabrik * 100

	// Sampai Hari Ini - Kering
	TotalSampaiHariIniK3Sheet       float64 `json:"totalSampaiHariIniK3Sheet"` // SUM dari K3Sheet (bukan perhitungan)
	TotalSampaiHariIniKeringSheet   float64 `json:"totalSampaiHariIniKeringSheet"`
	TotalSampaiHariIniKeringBrCr    float64 `json:"totalSampaiHariIniKeringBrCr"`
	TotalSampaiHariIniKeringJumlah  float64 `json:"totalSampaiHariIniKeringJumlah"`
	TotalSampaiHariIniK3SheetPersen float64 `json:"totalSampaiHariIniK3SheetPersen"` // KeringSheet / BasahLatekPabrik * 100

	// Produksi Per Taper
	TotalProduksiPerTaperHariIni       float64 `json:"totalProduksiPerTaperHariIni"`       // KeringJumlah / HKO
	TotalProduksiPerTaperSampaiHariIni float64 `json:"totalProduksiPerTaperSampaiHariIni"` // KeringJumlah / HKO

	TotalProduksiHariIni       float64 `json:"totalProduksiHariIni"`
	TotalProduksiSampaiHariIni float64 `json:"totalProduksiSampaiHariIni"`
}

// GetDashboardData mengambil data dashboard berdasarkan afdeling untuk tanggal hari ini
func GetDashboardData(w http.ResponseWriter, r *http.Request) {
	// Ambil parameter afdeling dari query string
	afdeling := r.URL.Query().Get("afdeling")
	if afdeling == "" {
		http.Error(w, "Parameter afdeling diperlukan", http.StatusBadRequest)
		return
	}

	// Debug logging
	log.Printf("📍 Request untuk afdeling: '%s'", afdeling)

	// Dapatkan tanggal hari ini (tanpa waktu)
	today := time.Now().Truncate(24 * time.Hour)
	log.Printf("📅 Tanggal query: %s", today.Format("2006-01-02"))

	db := config.GetDB()

	// Inisialisasi response
	response := dashboardDataResponse{}

	// Gunakan agregasi langsung di database untuk menghindari duplikasi
	type AggregateResult struct {
		TotalHKOHariIni                    int
		TotalHKOSampaiHariIni              int
		TotalHariIniBasahLatekKebun        float64
		TotalHariIniBasahLatekPabrik       float64
		TotalHariIniBasahLumpKebun         float64
		TotalHariIniBasahLumpPabrik        float64
		TotalHariIniK3Sheet                float64
		TotalHariIniKeringSheet            float64
		TotalHariIniKeringBrCr             float64
		TotalHariIniKeringJumlah           float64
		TotalProduksiHariIni               float64
		TotalSampaiHariIniBasahLatekKebun  float64
		TotalSampaiHariIniBasahLatekPabrik float64
		TotalSampaiHariIniBasahLumpKebun   float64
		TotalSampaiHariIniBasahLumpPabrik  float64
		TotalSampaiHariIniK3Sheet          float64
		TotalSampaiHariIniKeringSheet      float64
		TotalSampaiHariIniKeringBrCr       float64
		TotalSampaiHariIniKeringJumlah     float64
		TotalProduksiSampaiHariIni         float64
	}

	var result AggregateResult

	// FIX: Gunakan LOWER() untuk case-insensitive comparison
	err := db.Model(&models.Rekap{}).
		Select(`
			COALESCE(SUM(hko_hari_ini), 0) as total_hko_hari_ini,
			COALESCE(SUM(hko_sampai_hari_ini), 0) as total_hko_sampai_hari_ini,
			COALESCE(SUM(hari_ini_basah_latek_kebun), 0) as total_hari_ini_basah_latek_kebun,
			COALESCE(SUM(hari_ini_basah_latek_pabrik), 0) as total_hari_ini_basah_latek_pabrik,
			COALESCE(SUM(hari_ini_basah_lump_kebun), 0) as total_hari_ini_basah_lump_kebun,
			COALESCE(SUM(hari_ini_basah_lump_pabrik), 0) as total_hari_ini_basah_lump_pabrik,
			COALESCE(SUM(hari_ini_k3_sheet), 0) as total_hari_ini_k3_sheet,
			COALESCE(SUM(hari_ini_kering_sheet), 0) as total_hari_ini_kering_sheet,
			COALESCE(SUM(hari_ini_kering_br_cr), 0) as total_hari_ini_kering_br_cr,
			COALESCE(SUM(hari_ini_kering_jumlah), 0) as total_hari_ini_kering_jumlah,
			COALESCE(SUM(sampai_hari_ini_basah_latek_kebun), 0) as total_sampai_hari_ini_basah_latek_kebun,
			COALESCE(SUM(sampai_hari_ini_basah_latek_pabrik), 0) as total_sampai_hari_ini_basah_latek_pabrik,
			COALESCE(SUM(sampai_hari_ini_basah_lump_kebun), 0) as total_sampai_hari_ini_basah_lump_kebun,
			COALESCE(SUM(sampai_hari_ini_basah_lump_pabrik), 0) as total_sampai_hari_ini_basah_lump_pabrik,
			COALESCE(SUM(sampai_hari_ini_k3_sheet), 0) as total_sampai_hari_ini_k3_sheet,
			COALESCE(SUM(sampai_hari_ini_kering_sheet), 0) as total_sampai_hari_ini_kering_sheet,
			COALESCE(SUM(sampai_hari_ini_kering_br_cr), 0) as total_sampai_hari_ini_kering_br_cr,
			COALESCE(SUM(sampai_hari_ini_kering_jumlah), 0) as total_sampai_hari_ini_kering_jumlah,
			COALESCE(SUM(total_produksi_hari_ini),0) as total_produksi_hari_ini,
			COALESCE(SUM(total_produksi_sampai_hari_ini),0) as total_produksi_sampai_hari_ini
		`).
		Where("DATE(tanggal) = DATE(?) AND LOWER(afdeling) = LOWER(?) AND tipe_produksi != ?", today, afdeling, "REKAPITULASI").
		Scan(&result).Error

	if err != nil {
		log.Printf("❌ Error query database: %v", err)
		http.Error(w, "Error mengambil data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Debug logging hasil query
	log.Printf("📊 Data ditemukan - HKO Hari Ini: %d, Basah Latek Kebun: %.2f",
		result.TotalHKOHariIni, result.TotalHariIniBasahLatekKebun)

	// Copy hasil agregasi ke response
	response.TotalHKOHariIni = result.TotalHKOHariIni
	response.TotalHKOSampaiHariIni = result.TotalHKOSampaiHariIni
	response.TotalHariIniBasahLatekKebun = result.TotalHariIniBasahLatekKebun
	response.TotalHariIniBasahLatekPabrik = result.TotalHariIniBasahLatekPabrik
	response.TotalHariIniBasahLumpKebun = result.TotalHariIniBasahLumpKebun
	response.TotalHariIniBasahLumpPabrik = result.TotalHariIniBasahLumpPabrik
	response.TotalHariIniK3Sheet = result.TotalHariIniK3Sheet
	response.TotalHariIniKeringSheet = result.TotalHariIniKeringSheet
	response.TotalHariIniKeringBrCr = result.TotalHariIniKeringBrCr
	response.TotalHariIniKeringJumlah = result.TotalHariIniKeringJumlah
	response.TotalSampaiHariIniBasahLatekKebun = result.TotalSampaiHariIniBasahLatekKebun
	response.TotalSampaiHariIniBasahLatekPabrik = result.TotalSampaiHariIniBasahLatekPabrik
	response.TotalSampaiHariIniBasahLumpKebun = result.TotalSampaiHariIniBasahLumpKebun
	response.TotalSampaiHariIniBasahLumpPabrik = result.TotalSampaiHariIniBasahLumpPabrik
	response.TotalSampaiHariIniK3Sheet = result.TotalSampaiHariIniK3Sheet
	response.TotalSampaiHariIniKeringSheet = result.TotalSampaiHariIniKeringSheet
	response.TotalSampaiHariIniKeringBrCr = result.TotalSampaiHariIniKeringBrCr
	response.TotalSampaiHariIniKeringJumlah = result.TotalSampaiHariIniKeringJumlah

	// ✅ TAMBAHKAN 2 BARIS INI
	response.TotalProduksiHariIni = result.TotalProduksiHariIni
	response.TotalProduksiSampaiHariIni = result.TotalProduksiSampaiHariIni

	// === HITUNG PERSENTASE HARI INI ===

	// Persentase Basah Latek Hari Ini: (Kebun - Pabrik) / kebun * 100
	if response.TotalHariIniBasahLatekKebun > 0 {
		selisih := response.TotalHariIniBasahLatekKebun - response.TotalHariIniBasahLatekPabrik
		response.TotalHariIniBasahLatekPersen = (selisih / response.TotalHariIniBasahLatekKebun) * 100
	}

	// Persentase Basah Lump Hari Ini: (Kebun - Pabrik) / kebun * 100
	if response.TotalHariIniBasahLumpKebun > 0 {
		selisih := response.TotalHariIniBasahLumpKebun - response.TotalHariIniBasahLumpPabrik
		response.TotalHariIniBasahLumpPersen = (selisih / response.TotalHariIniBasahLumpKebun) * 100
	}

	// K3 Sheet Hari Ini: KeringSheet / BasahLatekPabrik * 100
	if response.TotalHariIniBasahLatekPabrik > 0 {
		response.TotalHariIniK3SheetPersen = (response.TotalHariIniKeringSheet / response.TotalHariIniBasahLatekPabrik) * 100
	}

	// Produksi Per Taper Hari Ini: KeringJumlah / HKO
	if response.TotalHKOHariIni > 0 {
		response.TotalProduksiPerTaperHariIni = response.TotalHariIniKeringJumlah / float64(response.TotalHKOHariIni)
	}

	// === HITUNG PERSENTASE SAMPAI HARI INI ===

	// Persentase Basah Latek Sampai Hari Ini: (Kebun - Pabrik) / Kebun * 100
	if response.TotalSampaiHariIniBasahLatekKebun > 0 {
		selisih := response.TotalSampaiHariIniBasahLatekKebun - response.TotalSampaiHariIniBasahLatekPabrik
		response.TotalSampaiHariIniBasahLatekPersen = (selisih / response.TotalSampaiHariIniBasahLatekKebun) * 100
	}

	// Persentase Basah Lump Sampai Hari Ini: (Kebun - Pabrik) / Kebun * 100
	if response.TotalSampaiHariIniBasahLumpKebun > 0 {
		selisih := response.TotalSampaiHariIniBasahLumpKebun - response.TotalSampaiHariIniBasahLumpPabrik
		response.TotalSampaiHariIniBasahLumpPersen = (selisih / response.TotalSampaiHariIniBasahLumpKebun) * 100
	}

	// K3 Sheet Sampai Hari Ini: KeringSheet / BasahLatekPabrik * 100
	if response.TotalSampaiHariIniBasahLatekPabrik > 0 {
		response.TotalSampaiHariIniK3SheetPersen = (response.TotalSampaiHariIniKeringSheet / response.TotalSampaiHariIniBasahLatekPabrik) * 100
	}

	// Produksi Per Taper Sampai Hari Ini: KeringJumlah / HKO
	if response.TotalHKOSampaiHariIni > 0 {
		response.TotalProduksiPerTaperSampaiHariIni = response.TotalSampaiHariIniKeringJumlah / float64(response.TotalHKOSampaiHariIni)
	}

	// Debug logging response
	log.Printf("✅ Response untuk %s: Total records dengan data non-zero", afdeling)

	// Kirim response sebagai JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
