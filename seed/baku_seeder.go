package seed

import (
	"app-inputan-ptpn/config"
	"app-inputan-ptpn/controllers"
	"bytes"
	"encoding/json"
	"net/http/httptest"
)

func SeedBaku() {
	config.DB.Exec("DELETE FROM baku_penyadaps")
	config.DB.Exec("DELETE FROM baku_details")
	entries := []map[string]interface{}{
		{"IdBakuMandor": 2, "IdPenyadap": 1, "BasahLatex": 39.40, "Sheet": 11.34, "BasahLump": 0, "BrCr": 0},
		{"IdBakuMandor": 2, "IdPenyadap": 2, "BasahLatex": 32.51, "Sheet": 9.36, "BasahLump": 1.75, "BrCr": 0.50},
		{"IdBakuMandor": 2, "IdPenyadap": 3, "BasahLatex": 28.57, "Sheet": 8.23, "BasahLump": 3.50, "BrCr": 1},
		{"IdBakuMandor": 2, "IdPenyadap": 4, "BasahLatex": 31.52, "Sheet": 9.07, "BasahLump": 1.75, "BrCr": 0.50},

		{"IdBakuMandor": 3, "IdPenyadap": 5, "BasahLatex": 37.50, "Sheet": 10.60, "BasahLump": 5.37, "BrCr": 1.90},
		{"IdBakuMandor": 3, "IdPenyadap": 6, "BasahLatex": 30.60, "Sheet": 8.66, "BasahLump": 4.47, "BrCr": 1.58},
		{"IdBakuMandor": 3, "IdPenyadap": 7, "BasahLatex": 44.42, "Sheet": 12.57, "BasahLump": 7.16, "BrCr": 2.52},
		{"IdBakuMandor": 3, "IdPenyadap": 8, "BasahLatex": 39.48, "Sheet": 11.17, "BasahLump": 0, "BrCr": 0},

		{"IdBakuMandor": 4, "IdPenyadap": 11, "BasahLatex": 9.30, "Sheet": 2.09, "BasahLump": 3.38, "BrCr": 1.23},
		{"IdBakuMandor": 4, "IdPenyadap": 12, "BasahLatex": 16.74, "Sheet": 3.77, "BasahLump": 5.92, "BrCr": 2.15},
		{"IdBakuMandor": 4, "IdPenyadap": 14, "BasahLatex": 13.96, "Sheet": 3.14, "BasahLump": 12.70, "BrCr": 4.62},

		{"IdBakuMandor": 7, "IdPenyadap": 21, "BasahLatex": 31.0, "Sheet": 7, "BasahLump": 10, "BrCr": 3},

		{"IdBakuMandor": 8, "IdPenyadap": 22, "BasahLatex": 39.41, "Sheet": 8.82, "BasahLump": 0, "BrCr": 0},
		{"IdBakuMandor": 8, "IdPenyadap": 24, "BasahLatex": 27.59, "Sheet": 6.18, "BasahLump": 5, "BrCr": 2},

		{"IdBakuMandor": 9, "IdPenyadap": 26, "BasahLatex": 25, "Sheet": 5.52, "BasahLump": 6, "BrCr": 2.1},
		{"IdBakuMandor": 9, "IdPenyadap": 27, "BasahLatex": 75, "Sheet": 16.56, "BasahLump": 2, "BrCr": 0.7},
		{"IdBakuMandor": 9, "IdPenyadap": 28, "BasahLatex": 41, "Sheet": 9.05, "BasahLump": 5, "BrCr": 1.74},
		{"IdBakuMandor": 9, "IdPenyadap": 29, "BasahLatex": 52, "Sheet": 17.48, "BasahLump": 8, "BrCr": 2.79},
		{"IdBakuMandor": 9, "IdPenyadap": 30, "BasahLatex": 47, "Sheet": 10.38, "BasahLump": 4, "BrCr": 1.4},
		{"IdBakuMandor": 9, "IdPenyadap": 31, "BasahLatex": 18, "Sheet": 3.97, "BasahLump": 12, "BrCr": 4.19},
		{"IdBakuMandor": 9, "IdPenyadap": 32, "BasahLatex": 17, "Sheet": 3.75, "BasahLump": 23, "BrCr": 8.03},
		{"IdBakuMandor": 9, "IdPenyadap": 33, "BasahLatex": 37, "Sheet": 9, "BasahLump": 2, "BrCr": 1},
		{"IdBakuMandor": 9, "IdPenyadap": 34, "BasahLatex": 33, "Sheet": 7.29, "BasahLump": 3, "BrCr": 1.05},

		{"IdBakuMandor": 10, "IdPenyadap": 35, "BasahLatex": 33.39, "Sheet": 7.97, "BasahLump": 0.8, "BrCr": 0.28},
		{"IdBakuMandor": 10, "IdPenyadap": 36, "BasahLatex": 58.92, "Sheet": 14.1, "BasahLump": 4.7, "BrCr": 1.62},
		{"IdBakuMandor": 10, "IdPenyadap": 37, "BasahLatex": 41.24, "Sheet": 9.87, "BasahLump": 4.7, "BrCr": 1.62},
		{"IdBakuMandor": 10, "IdPenyadap": 38, "BasahLatex": 39.28, "Sheet": 9.4, "BasahLump": 2.35, "BrCr": 0.81},
		{"IdBakuMandor": 10, "IdPenyadap": 39, "BasahLatex": 43.2, "Sheet": 10.34, "BasahLump": 3.92, "BrCr": 1.35},
		{"IdBakuMandor": 10, "IdPenyadap": 40, "BasahLatex": 54.01, "Sheet": 12.92, "BasahLump": 0, "BrCr": 0},
		{"IdBakuMandor": 10, "IdPenyadap": 41, "BasahLatex": 36.33, "Sheet": 8.69, "BasahLump": 4.7, "BrCr": 1.62},
		{"IdBakuMandor": 10, "IdPenyadap": 42, "BasahLatex": 37.31, "Sheet": 8.93, "BasahLump": 0.78, "BrCr": 0.27},
		{"IdBakuMandor": 10, "IdPenyadap": 44, "BasahLatex": 31.42, "Sheet": 7.52, "BasahLump": 2.35, "BrCr": 0.81},
		{"IdBakuMandor": 10, "IdPenyadap": 46, "BasahLatex": 41.24, "Sheet": 9.87, "BasahLump": 3.14, "BrCr": 1.08},
		{"IdBakuMandor": 10, "IdPenyadap": 47, "BasahLatex": 29.46, "Sheet": 7.05, "BasahLump": 0.78, "BrCr": 0.27},
		{"IdBakuMandor": 10, "IdPenyadap": 48, "BasahLatex": 43.2, "Sheet": 10.34, "BasahLump": 0.78, "BrCr": 0.27},

		{"IdBakuMandor": 11, "IdPenyadap": 52, "BasahLatex": 19, "Sheet": 4, "BasahLump": 8, "BrCr": 3},

		{"IdBakuMandor": 12, "IdPenyadap": 56, "BasahLatex": 25.31, "Sheet": 5.62, "BasahLump": 4.44, "BrCr": 1.67},
		{"IdBakuMandor": 12, "IdPenyadap": 57, "BasahLatex": 19.69, "Sheet": 4.38, "BasahLump": 3.56, "BrCr": 1.33},

		{"IdBakuMandor": 13, "IdPenyadap": 64, "BasahLatex": 30, "Sheet": 6, "BasahLump": 14, "BrCr": 5},
	}

	for _, e := range entries {
		body, _ := json.Marshal(e)
		req := httptest.NewRequest("POST", "/api/baku", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		controllers.CreateBakuPenyadap(w, req) // panggil fungsi controller langsung
		res := w.Result()
		defer res.Body.Close()
	}
}
