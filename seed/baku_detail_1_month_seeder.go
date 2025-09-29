package seed

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func SeedData() {
	fmt.Println("→ Memulai login ke API...")

	// 1. Login ke API dengan retry mechanism
	var loginSuccess bool
	var sessionToken string

	for attempt := 1; attempt <= 3; attempt++ {
		loginReq := map[string]string{
			"username": "admin",
			"password": "admin123",
		}
		loginBody, _ := json.Marshal(loginReq)

		resp, err := http.Post("http://localhost:8080/login", "application/json", bytes.NewBuffer(loginBody))
		if err != nil {
			log.Printf("  ⚠️  Attempt %d/3 - Login error: %v\n", attempt, err)
			time.Sleep(2 * time.Second)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			log.Printf("  ⚠️  Attempt %d/3 - Login failed, status: %d, body: %s\n", attempt, resp.StatusCode, string(bodyBytes))
			time.Sleep(2 * time.Second)
			continue
		}

		var loginResp struct {
			Success bool   `json:"success"`
			Token   string `json:"token"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
			log.Printf("  ⚠️  Attempt %d/3 - Decode error: %v\n", attempt, err)
			time.Sleep(2 * time.Second)
			continue
		}

		if !loginResp.Success {
			log.Printf("  ⚠️  Attempt %d/3 - Login response failed: %+v\n", attempt, loginResp)
			time.Sleep(2 * time.Second)
			continue
		}

		// Login berhasil
		sessionToken = loginResp.Token
		loginSuccess = true
		break
	}

	if !loginSuccess {
		log.Println("  ❌ Login gagal setelah 3 percobaan. Seeding dibatalkan.")
		return
	}

	// Simpan token global
	SessionToken = sessionToken
	fmt.Printf("  ✓ Login berhasil (Token: %s...)\n", sessionToken[:20])

	// 2. Validasi bahwa Mandor ID 1 dan Penyadap ID 1 ada
	fmt.Println("\n→ Validasi data master...")

	client := &http.Client{}

	// Cek Mandor ID 1 (dengan autentikasi)
	mandorReq, _ := http.NewRequest("GET", "http://localhost:8080/api/mandor/1", nil)
	mandorReq.AddCookie(&http.Cookie{Name: "session_token", Value: SessionToken})

	mandorResp, err := client.Do(mandorReq)
	if err != nil || mandorResp.StatusCode != http.StatusOK {
		log.Printf("  ❌ CRITICAL: Mandor ID 1 tidak ditemukan!\n")
		if err != nil {
			log.Printf("     Error: %v\n", err)
		} else {
			log.Printf("     Status: %d\n", mandorResp.StatusCode)
		}
		log.Println("     Solusi: Pastikan seed.SeedMandor() menghasilkan mandor dengan ID 1")
		return
	}
	mandorResp.Body.Close()
	fmt.Println("  ✓ Mandor ID 1 ditemukan")

	// Cek Penyadap ID 1
	penyadapReq, _ := http.NewRequest("GET", "http://localhost:8080/api/penyadap", nil)
	penyadapReq.AddCookie(&http.Cookie{Name: "session_token", Value: SessionToken})

	penyadapResp, err := client.Do(penyadapReq)
	if err != nil || penyadapResp.StatusCode != http.StatusOK {
		log.Printf("  ❌ CRITICAL: Tidak bisa akses data penyadap!\n")
		if err != nil {
			log.Printf("     Error: %v\n", err)
		} else {
			log.Printf("     Status: %d\n", penyadapResp.StatusCode)
		}
		return
	}
	penyadapResp.Body.Close()
	fmt.Println("  ✓ Data penyadap dapat diakses")

	// 3. Siapkan tanggal sebulan penuh
	now := time.Now()
	year, month := now.Year(), now.Month()
	startOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, now.Location())
	startOfNextMonth := startOfMonth.AddDate(0, 1, 0)
	endOfMonth := startOfNextMonth.AddDate(0, 0, -1)

	successCount := 0
	errorCount := 0
	totalDays := 0
	errorDetails := make(map[string]int) // Untuk track jenis error

	fmt.Printf("\n→ Mulai seeding data dari %s sampai %s\n",
		startOfMonth.Format("2006-01-02"),
		endOfMonth.Format("2006-01-02"))

	// 4. Loop tanggal & kirim data ke API
	for d := startOfMonth; !d.After(endOfMonth); d = d.AddDate(0, 0, 1) {
		totalDays++

		// PERBAIKAN: Pastikan format tanggal konsisten
		penyadap := map[string]interface{}{
			"IdBakuMandor": 43,
			"IdPenyadap":   1,
			"Tanggal":      d.Format("2006-01-02T15:04:05Z07:00"), // ISO 8601 format
			"Tipe":         "BAKU_INTERNAL",
			"TahunTanam":   2020,
			"BasahLatex":   10.0,
			"Sheet":        20.0,
			"BasahLump":    15.0,
			"BrCr":         5.0,
		}

		body, _ := json.Marshal(penyadap)
		req, _ := http.NewRequest("POST", "http://localhost:8080/api/baku", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{
			Name:  "session_token",
			Value: SessionToken,
		})

		res, err := client.Do(req)
		if err != nil {
			errorMsg := fmt.Sprintf("Network error: %v", err)
			fmt.Printf("  ✗ Tanggal %s: %s\n", d.Format("2006-01-02"), errorMsg)
			errorDetails[errorMsg]++
			errorCount++
			time.Sleep(500 * time.Millisecond) // Delay lebih lama jika error
			continue
		}

		if res.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(res.Body)
			errorMsg := fmt.Sprintf("Status %d: %s", res.StatusCode, string(bodyBytes))

			// Hanya tampilkan error pertama untuk setiap jenis
			if errorDetails[errorMsg] == 0 {
				fmt.Printf("  ✗ Tanggal %s: %s\n", d.Format("2006-01-02"), errorMsg)
			}
			errorDetails[errorMsg]++
			errorCount++
		} else {
			if totalDays%5 == 0 || totalDays == 1 { // Tampilkan di hari 1, 5, 10, dst
				fmt.Printf("  ✓ Progress: %d/%d hari berhasil\n", successCount+1, endOfMonth.Day())
			}
			successCount++
		}
		res.Body.Close()

		// Delay kecil untuk menghindari overload
		time.Sleep(50 * time.Millisecond)
	}

	// 5. Summary dengan detail error
	fmt.Printf("\n→ Seeding data selesai: %d berhasil, %d gagal dari %d hari\n",
		successCount, errorCount, totalDays)

	if errorCount > 0 {
		fmt.Println("\n⚠️  DETAIL ERROR:")
		for errorMsg, count := range errorDetails {
			fmt.Printf("   - %s (terjadi %d kali)\n", errorMsg, count)
		}

		fmt.Println("\n💡 KEMUNGKINAN PENYEBAB & SOLUSI:")
		fmt.Println("   1. Mandor ID 1 tidak ada → Jalankan seed.SeedMandor() dengan benar")
		fmt.Println("   2. Penyadap ID 1 tidak ada → Jalankan seed.SeedPenyadap() dengan benar")
		fmt.Println("   3. Tipe 'BAKU' tidak valid → Cek models/baku.go untuk TipeProduksi")
		fmt.Println("   4. Session expired → Coba perpanjang waktu tunggu di main.go")
		fmt.Println("   5. Duplikasi data → Cek apakah ada constraint unique yang dilanggar")
	}
}
