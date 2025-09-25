package seed

import (
	"app-inputan-ptpn/config"
	"app-inputan-ptpn/controllers"
	"app-inputan-ptpn/models"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
)

func SeedBakuBorong() {
	entries := []map[string]interface{}{
		{"IdBakuMandor": 98, "IdPenyadap": 13, "BasahLatex": 13.15, "Sheet": 2.97, "BasahLump": 6, "BrCr": 2, "Tipe": "BAKU_BORONG"},
		{"IdBakuMandor": 98, "IdPenyadap": 15, "BasahLatex": 8.45, "Sheet": 1.91, "BasahLump": 3.43, "BrCr": 1.14, "Tipe": "BAKU_BORONG"},
		{"IdBakuMandor": 98, "IdPenyadap": 16, "BasahLatex": 9.4, "Sheet": 2.12, "BasahLump": 2.57, "BrCr": 0.86, "Tipe": "BAKU_BORONG"},

		{"IdBakuMandor": 100, "IdPenyadap": 18, "BasahLatex": 27.74, "Sheet": 5.81, "BasahLump": 0, "BrCr": 0, "Tipe": "BAKU_BORONG"},
		{"IdBakuMandor": 100, "IdPenyadap": 19, "BasahLatex": 30.52, "Sheet": 6.38, "BasahLump": 0, "BrCr": 0, "Tipe": "BAKU_BORONG"},
		{"IdBakuMandor": 100, "IdPenyadap": 20, "BasahLatex": 27.74, "Sheet": 5.81, "BasahLump": 0, "BrCr": 0, "Tipe": "BAKU_BORONG"},

		{"IdBakuMandor": 104, "IdPenyadap": 45, "BasahLatex": 20, "Sheet": 5, "BasahLump": 2, "BrCr": 1, "Tipe": "BAKU_BORONG"},

		{"IdBakuMandor": 105, "IdPenyadap": 50, "BasahLatex": 13.96, "Sheet": 2.92, "BasahLump": 2.60, "BrCr": 0.91, "Tipe": "BAKU_BORONG"},
		{"IdBakuMandor": 105, "IdPenyadap": 51, "BasahLatex": 18.61, "Sheet": 3.89, "BasahLump": 8.7, "BrCr": 3.05, "Tipe": "BAKU_BORONG"},
		{"IdBakuMandor": 105, "IdPenyadap": 53, "BasahLatex": 16.75, "Sheet": 3.5, "BasahLump": 6.09, "BrCr": 2.13, "Tipe": "BAKU_BORONG"},
		{"IdBakuMandor": 105, "IdPenyadap": 54, "BasahLatex": 17.68, "Sheet": 3.69, "BasahLump": 2.61, "BrCr": 0.91, "Tipe": "BAKU_BORONG"},

		{"IdBakuMandor": 26, "IdPenyadap": 60, "BasahLatex": 16.88, "Sheet": 3.56, "BasahLump": 4.17, "BrCr": 1.46, "Tipe": "BAKU_BORONG"},
		{"IdBakuMandor": 26, "IdPenyadap": 61, "BasahLatex": 18.75, "Sheet": 3.96, "BasahLump": 4.17, "BrCr": 1.46, "Tipe": "BAKU_BORONG"},
		{"IdBakuMandor": 26, "IdPenyadap": 63, "BasahLatex": 21.56, "Sheet": 4.55, "BasahLump": 1.66, "BrCr": 0.58, "Tipe": "BAKU_BORONG"},
		{"IdBakuMandor": 26, "IdPenyadap": 65, "BasahLatex": 21.56, "Sheet": 4.55, "BasahLump": 5, "BrCr": 1.75, "Tipe": "BAKU_BORONG"},
		{"IdBakuMandor": 26, "IdPenyadap": 66, "BasahLatex": 11.25, "Sheet": 2.38, "BasahLump": 5, "BrCr": 1.75, "Tipe": "BAKU_BORONG"},
	}

	for _, e := range entries {
		// Periksa apakah data dengan kombinasi IdBakuMandor, IdPenyadap, dan Tipe sudah ada
		var existing models.BakuPenyadap
		config.DB.Where("id_baku_mandor = ? AND id_penyadap = ? AND tipe = ?", e["IdBakuMandor"], e["IdPenyadap"], e["Tipe"]).First(&existing)

		// Jika data belum ada, lanjutkan untuk menyimpan
		if existing.ID == 0 {
			body, _ := json.Marshal(e)
			req := httptest.NewRequest("POST", "/api/baku", bytes.NewBuffer(body))
			w := httptest.NewRecorder()
			controllers.CreateBakuPenyadap(w, req) // Memanggil fungsi controller langsung
			res := w.Result()
			defer res.Body.Close()
		} else {
			// Data sudah ada, bisa menambahkan logika lain jika perlu, misalnya update atau log
			fmt.Printf("Data dengan IdBakuMandor %v, IdPenyadap %v, dan Tipe %v sudah ada. Skip.\n", e["IdBakuMandor"], e["IdPenyadap"], e["Tipe"])
		}
	}
}
