package seed

import (
	"app-inputan-ptpn/models"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func SeedData() {
	fmt.Println("=== MULAI SEEDING VIA API (dengan login) ===")

	// 1. Login ke API
	loginReq := map[string]string{
		"username": "admin",
		"password": "admin123",
	}
	loginBody, _ := json.Marshal(loginReq)

	resp, err := http.Post("http://localhost:8080/login", "application/json", bytes.NewBuffer(loginBody))
	if err != nil {
		log.Fatalf("Login gagal: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Login gagal, status: %d", resp.StatusCode)
	}

	var loginResp struct {
		Success bool   `json:"success"`
		Token   string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		log.Fatalf("Gagal decode response login: %v", err)
	}

	if !loginResp.Success {
		log.Fatalf("Login gagal, response: %+v", loginResp)
	}

	// Simpan token global
	SessionToken = loginResp.Token
	log.Println("✓ Login berhasil, token:", SessionToken)

	// 2. Siapkan tanggal sebulan penuh
	now := time.Now()
	year, month := now.Year(), now.Month()
	startOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, now.Location())
	startOfNextMonth := startOfMonth.AddDate(0, 1, 0)
	endOfMonth := startOfNextMonth.AddDate(0, 0, -1)

	client := &http.Client{}
	successCount := 0
	errorCount := 0

	// 3. Loop tanggal & kirim data ke API
	for d := startOfMonth; !d.After(endOfMonth); d = d.AddDate(0, 0, 1) {
		penyadap := models.BakuPenyadap{
			IdBakuMandor: 1,
			IdPenyadap:   1,
			Tanggal:      d,
			Tipe:         models.TipeBaku,
			TahunTanam:   2020,
			BasahLatex:   10.0,
			Sheet:        20.0,
			BasahLump:    15.0,
			BrCr:         5.0,
		}

		body, _ := json.Marshal(penyadap)
		req, _ := http.NewRequest("POST", "http://localhost:8080/api/baku", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		// pakai cookie dari token
		req.AddCookie(&http.Cookie{
			Name:  "session_token",
			Value: SessionToken,
		})

		res, err := client.Do(req)
		if err != nil {
			fmt.Printf("✗ Error kirim tanggal %s: %v\n", d.Format("2006-01-02"), err)
			errorCount++
			continue
		}
		if res.StatusCode != http.StatusCreated {
			fmt.Printf("✗ Gagal insert via API tanggal %s (status %d)\n", d.Format("2006-01-02"), res.StatusCode)
			errorCount++
		} else {
			fmt.Printf("✓ Berhasil insert via API tanggal %s\n", d.Format("2006-01-02"))
			successCount++
		}
		res.Body.Close()
	}

	// 4. Ringkasan
	fmt.Println("=== SELESAI SEEDING ===")
	fmt.Printf("Berhasil: %d\n", successCount)
	fmt.Printf("Gagal: %d\n", errorCount)
}
