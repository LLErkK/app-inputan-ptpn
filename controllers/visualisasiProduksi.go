package controllers

import (
	"app-inputan-ptpn/config"
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"time"

	"app-inputan-ptpn/models"
)

type VisualisasiProduksiResponse struct {
	Labels []string            `json:"labels"`
	Data   []ProduksiDataPoint `json:"data"`
}

type ProduksiDataPoint struct {
	Tanggal string  `json:"tanggal"`
	Value   float64 `json:"value"`
}

func getNikPenyadapById(idPenyadap uint) string {
	var penyadap models.Penyadap
	if err := config.DB.First(&penyadap, idPenyadap).Error; err != nil {
		return ""
	}
	return penyadap.NIK

}
func GetVisualisasiProduksi(w http.ResponseWriter, r *http.Request) {
	tanggalAwal := r.URL.Query().Get("tanggalAwal")
	tanggalAkhir := r.URL.Query().Get("tanggalAkhir")
	satuan := r.URL.Query().Get("satuan")
	tipeProduksi := r.URL.Query().Get("tipeProduksi")
	idPenyadap := r.URL.Query().Get("idPenyadap")

	uintIdPenyadap64, err := strconv.ParseUint(idPenyadap, 10, 64)
	if err != nil {
		http.Error(w, "Parameter idPenyadap tidak valid", http.StatusBadRequest)
		return
	}

	nikPenyadap := getNikPenyadapById(uint(uintIdPenyadap64))

	// Validasi parameter wajib
	if nikPenyadap == "" {
		http.Error(w, "Parameter nikPenyadap tidak boleh kosong", http.StatusBadRequest)
		return
	}
	if tanggalAwal == "" || tanggalAkhir == "" {
		http.Error(w, "Parameter tanggalAwal dan tanggalAkhir tidak boleh kosong", http.StatusBadRequest)
		return
	}
	if satuan == "" {
		http.Error(w, "Parameter satuan tidak boleh kosong", http.StatusBadRequest)
		return
	}

	// Validasi satuan - field yang tersedia di model Produksi
	validSatuan := map[string]bool{
		"basah_latek": true,
		"sheet":       true,
		"basah_lump":  true,
		"br_cr":       true,
	}

	if !validSatuan[satuan] {
		http.Error(w, "Parameter satuan tidak valid. Gunakan: basah_latek, sheet, basah_lump, atau br_cr", http.StatusBadRequest)
		return
	}

	result, err := visualisasiProduksiPenyadap(nikPenyadap, tipeProduksi, tanggalAwal, tanggalAkhir, satuan)
	if err != nil {
		http.Error(w, "Error mengambil data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func visualisasiProduksiPenyadap(nikPenyadap, tipeProduksi, tanggalAwal, tanggalAkhir, satuan string) (VisualisasiProduksiResponse, error) {
	var produksiList []models.Produksi
	db := config.GetDB()
	query := db.Model(&models.Produksi{})

	// Filter berdasarkan tanggal
	startDate, _ := time.Parse("2006-01-02", tanggalAwal)
	endDate, _ := time.Parse("2006-01-02", tanggalAkhir)
	startDate = startDate.AddDate(0, 0, -1)
	query = query.Where("tanggal BETWEEN ? AND ?", startDate, endDate)

	// Filter NIK Penyadap
	query = query.Where("nik = ?", nikPenyadap)

	// Filter tipeProduksi jika ada dan tidak "-"
	if tipeProduksi != "" && tipeProduksi != "-" {
		query = query.Where("tipe_produksi = ?", tipeProduksi)
	}

	if err := query.Order("tanggal ASC").Find(&produksiList).Error; err != nil {
		return VisualisasiProduksiResponse{}, err
	}
	if len(produksiList) == 0 {
		return VisualisasiProduksiResponse{
			Labels: []string{},
			Data:   []ProduksiDataPoint{},
		}, nil
	}

	return aggregateProduksiData(produksiList, satuan), nil
}

func aggregateProduksiData(produksiList []models.Produksi, satuan string) VisualisasiProduksiResponse {
	dataMap := make(map[string]float64)

	// Agregasi data per tanggal
	for _, produksi := range produksiList {
		dateStr := produksi.Tanggal.Format("2006-01-02")

		switch satuan {
		case "basah_latek":
			dataMap[dateStr] += produksi.BasahLatek
		case "sheet":
			dataMap[dateStr] += produksi.Sheet
		case "basah_lump":
			dataMap[dateStr] += produksi.BasahLump
		case "br_cr":
			dataMap[dateStr] += produksi.BrCr
		}
	}

	// Extract dan sort tanggal
	var dates []string
	for date := range dataMap {
		dates = append(dates, date)
	}

	// Sort tanggal secara ascending
	sort.Strings(dates)

	// Build data berdasarkan urutan tanggal
	var data []ProduksiDataPoint
	for _, date := range dates {
		data = append(data, ProduksiDataPoint{
			Tanggal: date,
			Value:   dataMap[date],
		})
	}

	return VisualisasiProduksiResponse{
		Labels: dates,
		Data:   data,
	}
}
