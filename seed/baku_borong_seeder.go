package seed

import (
	"app-inputan-ptpn/controllers"
	"bytes"
	"encoding/json"
	"net/http/httptest"
)

func SeedBakuBorong() {
	entries := []map[string]interface{}{
		{"IdBakuMandor": 98, "IdPenyadap": 13, "BasahLatex": 13.15, "Sheet": 2.97, "BasahLump": 6, "BrCr": 2},
		{"IdBakuMandor": 98, "IdPenyadap": 15, "BasahLatex": 8.45, "Sheet": 1.91, "BasahLump": 3.43, "BrCr": 1.14},
		{"IdBakuMandor": 98, "IdPenyadap": 16, "BasahLatex": 9.4, "Sheet": 2.12, "BasahLump": 2.57, "BrCr": 0.86},

		{"IdBakuMandor": 100, "IdPenyadap": 18, "BasahLatex": 27.74, "Sheet": 5.81, "BasahLump": 0, "BrCr": 0},
		{"IdBakuMandor": 100, "IdPenyadap": 19, "BasahLatex": 30.52, "Sheet": 6.38, "BasahLump": 0, "BrCr": 0},
		{"IdBakuMandor": 100, "IdPenyadap": 20, "BasahLatex": 27.74, "Sheet": 5.81, "BasahLump": 0, "BrCr": 0},

		{"IdBakuMandor": 104, "IdPenyadap": 45, "BasahLatex": 20, "Sheet": 5, "BasahLump": 2, "BrCr": 1},

		{"IdBakuMandor": 105, "IdPenyadap": 50, "BasahLatex": 13.96, "Sheet": 2.92, "BasahLump": 2.60, "BrCr": 0.91},
		{"IdBakuMandor": 105, "IdPenyadap": 51, "BasahLatex": 18.61, "Sheet": 3.89, "BasahLump": 8.7, "BrCr": 3.05},
		{"IdBakuMandor": 105, "IdPenyadap": 53, "BasahLatex": 16.75, "Sheet": 3.5, "BasahLump": 6.09, "BrCr": 2.13},
		{"IdBakuMandor": 105, "IdPenyadap": 54, "BasahLatex": 17.68, "Sheet": 3.69, "BasahLump": 2.61, "BrCr": 0.91},

		{"IdBakuMandor": 26, "IdPenyadap": 60, "BasahLatex": 16.88, "Sheet": 3.56, "BasahLump": 4.17, "BrCr": 1.46},
		{"IdBakuMandor": 26, "IdPenyadap": 61, "BasahLatex": 18.75, "Sheet": 3.96, "BasahLump": 4.17, "BrCr": 1.46},
		{"IdBakuMandor": 26, "IdPenyadap": 63, "BasahLatex": 21.56, "Sheet": 4.55, "BasahLump": 1.66, "BrCr": 0.58},
		{"IdBakuMandor": 26, "IdPenyadap": 65, "BasahLatex": 21.56, "Sheet": 4.55, "BasahLump": 5, "BrCr": 1.75},
		{"IdBakuMandor": 26, "IdPenyadap": 66, "BasahLatex": 11.25, "Sheet": 2.38, "BasahLump": 5, "BrCr": 1.75},
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
