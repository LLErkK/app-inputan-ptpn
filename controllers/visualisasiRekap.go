package controllers

import (
	"app-inputan-ptpn/config"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"app-inputan-ptpn/models"
)

type VisualisasiResponse struct {
	Labels []string    `json:"labels"`
	Data   []DataPoint `json:"data"`
}

type DataPoint struct {
	Tanggal string  `json:"tanggal"`
	Value   float64 `json:"value"`
}

func getMandorByID(idMandor int) (string, string) {
	var mandor models.BakuMandor
	if err := config.DB.First(&mandor, idMandor).Error; err != nil {
		return "", ""
	}
	return mandor.NIK, strconv.Itoa(int(mandor.TahunTanam))
}

func GetVisualisasiRekap(w http.ResponseWriter, r *http.Request) {
	tipeData := r.URL.Query().Get("tipeData")
	tipeProduksi := r.URL.Query().Get("tipeProduksi")
	afdeling := r.URL.Query().Get("afdeling")
	idMandor := r.URL.Query().Get("idMandor")
	tanggalAwal := r.URL.Query().Get("tanggalAwal")
	tanggalAkhir := r.URL.Query().Get("tanggalAkhir")
	satuan := r.URL.Query().Get("satuan")
	//get nik dan tahuntanam
	intidMandor, _ := strconv.Atoi(idMandor)
	nikMandor, tahunTanam := getMandorByID(intidMandor)

	// Validasi parameter wajib
	if tipeData == "" {
		http.Error(w, "Parameter tipeData tidak boleh kosong", http.StatusBadRequest)
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

	// Validasi satuan - hanya field hari_ini yang valid
	validSatuan := map[string]bool{
		"hko":                true,
		"basah_latek_kebun":  true,
		"basah_latek_pabrik": true,
		"basah_latek_persen": true,
		"basah_lump_kebun":   true,
		"basah_lump_pabrik":  true,
		"basah_lump_persen":  true,
		"k3_sheet":           true,
		"kering_sheet":       true,
		"kering_br_cr":       true,
		"kering_jumlah":      true,
		"produksi_per_taper": true,
	}

	if !validSatuan[satuan] {
		http.Error(w, "Parameter satuan tidak valid", http.StatusBadRequest)
		return
	}

	var result VisualisasiResponse
	var err error

	switch tipeData {
	case "total":
		result, err = visualisasiTotal(tipeProduksi, tanggalAwal, tanggalAkhir, satuan)
	case "afdeling":
		if afdeling == "" {
			http.Error(w, "Parameter afdeling tidak boleh kosong untuk tipe 'afdeling'", http.StatusBadRequest)
			return
		}
		result, err = visualisasiAfdeling(tipeProduksi, afdeling, tanggalAwal, tanggalAkhir, satuan)
	case "mandor":
		if nikMandor == "" {
			http.Error(w, "Parameter nikMandor tidak boleh kosong untuk tipe 'mandor'", http.StatusBadRequest)
			return
		}
		result, err = visualisasiMandor(tipeProduksi, afdeling, nikMandor, tahunTanam, tanggalAwal, tanggalAkhir, satuan)
	default:
		http.Error(w, "Parameter tipeData tidak valid. Gunakan: total, afdeling, atau mandor", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "Error mengambil data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func visualisasiTotal(tipeProduksi, tanggalAwal, tanggalAkhir, satuan string) (VisualisasiResponse, error) {
	var rekaps []models.Rekap
	db := config.GetDB()
	query := db.Model(&models.Rekap{})

	startDate, _ := time.Parse("2006-01-02", tanggalAwal)
	endDate, _ := time.Parse("2006-01-02", tanggalAkhir)
	query = query.Where("tanggal BETWEEN ? AND ?", startDate, endDate)

	if tipeProduksi != "" && tipeProduksi != "-" {
		query = query.Where("tipe_produksi = ?", tipeProduksi)
	}

	if err := query.Order("tanggal ASC").Find(&rekaps).Error; err != nil {
		return VisualisasiResponse{}, err
	}

	return aggregateData(rekaps, satuan), nil
}

func visualisasiAfdeling(tipeProduksi, afdeling, tanggalAwal, tanggalAkhir, satuan string) (VisualisasiResponse, error) {
	var rekaps []models.Rekap
	db := config.GetDB()
	query := db.Model(&models.Rekap{})

	startDate, _ := time.Parse("2006-01-02", tanggalAwal)
	endDate, _ := time.Parse("2006-01-02", tanggalAkhir)
	query = query.Where("tanggal BETWEEN ? AND ?", startDate, endDate)
	query = query.Where("afdeling = ?", afdeling)

	if tipeProduksi != "" && tipeProduksi != "-" {
		query = query.Where("tipe_produksi = ?", tipeProduksi)
	}

	if err := query.Order("tanggal ASC").Find(&rekaps).Error; err != nil {
		return VisualisasiResponse{}, err
	}

	return aggregateData(rekaps, satuan), nil
}

func visualisasiMandor(tipeProduksi, afdeling, nikMandor, tahunTanam, tanggalAwal, tanggalAkhir, satuan string) (VisualisasiResponse, error) {
	var rekaps []models.Rekap
	db := config.GetDB()
	query := db.Model(&models.Rekap{})

	startDate, _ := time.Parse("2006-01-02", tanggalAwal)
	endDate, _ := time.Parse("2006-01-02", tanggalAkhir)
	query = query.Where("tanggal BETWEEN ? AND ?", startDate, endDate)
	query = query.Where("nik = ? AND tahun_tanam = ?", nikMandor, tahunTanam)

	if afdeling != "" && afdeling != "-" {
		query = query.Where("afdeling = ?", afdeling)
	}

	if tipeProduksi != "" && tipeProduksi != "-" {
		query = query.Where("tipe_produksi = ?", tipeProduksi)
	}

	if err := query.Order("tanggal ASC").Find(&rekaps).Error; err != nil {
		return VisualisasiResponse{}, err
	}

	return aggregateData(rekaps, satuan), nil
}

type tempDataPoint struct {
	HKO              int
	BasahLatekKebun  float64
	BasahLatekPabrik float64
	BasahLumpKebun   float64
	BasahLumpPabrik  float64
	KeringSheet      float64
	KeringBrCr       float64
	KeringJumlah     float64
}

func aggregateData(rekaps []models.Rekap, satuan string) VisualisasiResponse {
	dataMap := make(map[string]*tempDataPoint)

	// Agregasi data per tanggal
	for _, rekap := range rekaps {
		dateStr := rekap.Tanggal.Format("2006-01-02")

		if _, exists := dataMap[dateStr]; !exists {
			dataMap[dateStr] = &tempDataPoint{}
		}

		point := dataMap[dateStr]

		// Agregasi hanya data hari_ini
		point.HKO += rekap.HKOHariIni
		point.BasahLatekKebun += rekap.HariIniBasahLatekKebun
		point.BasahLatekPabrik += rekap.HariIniBasahLatekPabrik
		point.BasahLumpKebun += rekap.HariIniBasahLumpKebun
		point.BasahLumpPabrik += rekap.HariIniBasahLumpPabrik
		point.KeringSheet += rekap.HariIniKeringSheet
		point.KeringBrCr += rekap.HariIniKeringBrCr
		point.KeringJumlah += rekap.HariIniKeringJumlah
	}

	// Extract field yang diminta berdasarkan satuan
	var dates []string
	var data []DataPoint

	for date, point := range dataMap {
		dates = append(dates, date)

		var value float64

		switch satuan {
		case "hko":
			value = float64(point.HKO)
		case "basah_latek_kebun":
			value = point.BasahLatekKebun
		case "basah_latek_pabrik":
			value = point.BasahLatekPabrik
		case "basah_latek_persen":
			// Rumus: (Kebun - Pabrik) / Kebun * 100
			if point.BasahLatekKebun > 0 {
				value = ((point.BasahLatekKebun - point.BasahLatekPabrik) / point.BasahLatekKebun) * 100
			}
		case "basah_lump_kebun":
			value = point.BasahLumpKebun
		case "basah_lump_pabrik":
			value = point.BasahLumpPabrik
		case "basah_lump_persen":
			// Rumus: (Kebun - Pabrik) / Kebun * 100
			if point.BasahLumpKebun > 0 {
				value = ((point.BasahLumpKebun - point.BasahLumpPabrik) / point.BasahLumpKebun) * 100
			}
		case "k3_sheet":
			// Rumus: Kering Sheet / Basah Latek Pabrik * 100
			if point.BasahLatekPabrik > 0 {
				value = (point.KeringSheet / point.BasahLatekPabrik) * 100
			}
		case "kering_sheet":
			value = point.KeringSheet
		case "kering_br_cr":
			value = point.KeringBrCr
		case "kering_jumlah":
			value = point.KeringJumlah
		case "produksi_per_taper":
			// Rumus: Kering Jumlah / HKO
			if point.HKO > 0 {
				value = point.KeringJumlah / float64(point.HKO)
			}
		}

		data = append(data, DataPoint{
			Tanggal: date,
			Value:   value,
		})
	}

	return VisualisasiResponse{
		Labels: dates,
		Data:   data,
	}
}
