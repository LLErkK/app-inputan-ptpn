package seed

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func SeedBakuBorong() {
	fmt.Println("=== MULAI SEEDING BAKU BORONG VIA API ===")

	if SessionToken == "" {
		fmt.Println("✗ Token kosong, jalankan SeedData dulu")
		return
	}

	// Data entries
	entries := []map[string]interface{}{
		{"IdBakuMandor": 18, "IdPenyadap": 13, "BasahLatex": 13.15, "Sheet": 2.97, "BasahLump": 6, "BrCr": 2, "Tipe": "BAKU_BORONG"},
		{"IdBakuMandor": 18, "IdPenyadap": 15, "BasahLatex": 8.45, "Sheet": 1.91, "BasahLump": 3.43, "BrCr": 1.14, "Tipe": "BAKU_BORONG"},
		{"IdBakuMandor": 18, "IdPenyadap": 16, "BasahLatex": 9.4, "Sheet": 2.12, "BasahLump": 2.57, "BrCr": 0.86, "Tipe": "BAKU_BORONG"},
		{"IdBakuMandor": 19, "IdPenyadap": 18, "BasahLatex": 27.74, "Sheet": 5.81, "BasahLump": 0, "BrCr": 0, "Tipe": "BAKU_BORONG"},
		{"IdBakuMandor": 19, "IdPenyadap": 19, "BasahLatex": 30.52, "Sheet": 6.38, "BasahLump": 0, "BrCr": 0, "Tipe": "BAKU_BORONG"},
		{"IdBakuMandor": 19, "IdPenyadap": 20, "BasahLatex": 27.74, "Sheet": 5.81, "BasahLump": 0, "BrCr": 0, "Tipe": "BAKU_BORONG"},
		{"IdBakuMandor": 21, "IdPenyadap": 45, "BasahLatex": 20, "Sheet": 5, "BasahLump": 2, "BrCr": 1, "Tipe": "BAKU_BORONG"},
		{"IdBakuMandor": 24, "IdPenyadap": 50, "BasahLatex": 13.96, "Sheet": 2.92, "BasahLump": 2.60, "BrCr": 0.91, "Tipe": "BAKU_BORONG"},
		{"IdBakuMandor": 24, "IdPenyadap": 51, "BasahLatex": 18.61, "Sheet": 3.89, "BasahLump": 8.7, "BrCr": 3.05, "Tipe": "BAKU_BORONG"},
		{"IdBakuMandor": 24, "IdPenyadap": 53, "BasahLatex": 16.75, "Sheet": 3.5, "BasahLump": 6.09, "BrCr": 2.13, "Tipe": "BAKU_BORONG"},
		{"IdBakuMandor": 24, "IdPenyadap": 54, "BasahLatex": 17.68, "Sheet": 3.69, "BasahLump": 2.61, "BrCr": 0.91, "Tipe": "BAKU_BORONG"},
		{"IdBakuMandor": 26, "IdPenyadap": 60, "BasahLatex": 16.88, "Sheet": 3.56, "BasahLump": 4.17, "BrCr": 1.46, "Tipe": "BAKU_BORONG"},
		{"IdBakuMandor": 26, "IdPenyadap": 61, "BasahLatex": 18.75, "Sheet": 3.96, "BasahLump": 4.17, "BrCr": 1.46, "Tipe": "BAKU_BORONG"},
		{"IdBakuMandor": 26, "IdPenyadap": 63, "BasahLatex": 21.56, "Sheet": 4.55, "BasahLump": 1.66, "BrCr": 0.58, "Tipe": "BAKU_BORONG"},
		{"IdBakuMandor": 26, "IdPenyadap": 65, "BasahLatex": 21.56, "Sheet": 4.55, "BasahLump": 5, "BrCr": 1.75, "Tipe": "BAKU_BORONG"},
		{"IdBakuMandor": 26, "IdPenyadap": 66, "BasahLatex": 11.25, "Sheet": 2.38, "BasahLump": 5, "BrCr": 1.75, "Tipe": "BAKU_BORONG"},
	}

	client := &http.Client{}
	successCount := 0
	errorCount := 0

	// Loop insert ke API
	for _, e := range entries {
		body, _ := json.Marshal(e)
		req, _ := http.NewRequest("POST", "http://localhost:8080/api/baku", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{
			Name:  "session_token",
			Value: SessionToken,
		})

		res, err := client.Do(req)
		if err != nil {
			fmt.Printf("✗ Error insert BakuBorong (Mandor=%v, Penyadap=%v): %v\n", e["IdBakuMandor"], e["IdPenyadap"], err)
			errorCount++
			continue
		}
		if res.StatusCode != http.StatusCreated {
			fmt.Printf("✗ Gagal insert BakuBorong (Mandor=%v, Penyadap=%v) - status %d\n", e["IdBakuMandor"], e["IdPenyadap"], res.StatusCode)
			errorCount++
		} else {
			fmt.Printf("✓ Berhasil insert BakuBorong (Mandor=%v, Penyadap=%v)\n", e["IdBakuMandor"], e["IdPenyadap"])
			successCount++
		}
		res.Body.Close()
	}

	fmt.Println("=== SELESAI SEEDING BAKU BORONG ===")
	fmt.Printf("Berhasil: %d\n", successCount)
	fmt.Printf("Gagal: %d\n", errorCount)
}
