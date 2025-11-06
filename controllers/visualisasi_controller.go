package controllers

import (
	"net/http"
)

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func ServeVisualisasiPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/html/visualisasi.html")
}

// parameter tipeData, afdeling, namaMandor/namaPenyadap, tanggalAwal, tanggalAkhir
func GetVisualisasiData(w http.ResponseWriter, r *http.Request) {
	tipeData := r.URL.Query().Get("tipeData")

	if tipeData == "penyadap" {
		// TODO: ke controller visualisasi produksi
		GetVisualisasiProduksi(w, r)
	} else if tipeData == "mandor" || tipeData == "total" || tipeData == "afdeling" {
		// ke controller visualisasi rekap
		GetVisualisasiRekap(w, r)
	} else {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Parameter tipeData tidak valid. Gunakan: penyadap,mandor,total,afdeling",
		})
	}
}
