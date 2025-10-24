package seed

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
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

	SessionToken = sessionToken
	fmt.Printf("  ✓ Login berhasil (Token: %s...)\n", sessionToken[:20])

	// 2. Validasi master data
	fmt.Println("\n→ Validasi data master...")

	client := &http.Client{}

	// Cek Mandor ID 1
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
		return
	}
	mandorResp.Body.Close()
	fmt.Println("  ✓ Mandor ID 1 ditemukan")

	// Cek Penyadap
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

	// 3. Siapkan tanggal untuk 3 bulan (kemarin, ini, depan)
	now := time.Now()
	year, month, _ := now.Date()
	loc := now.Location()

	startOfPrevMonth := time.Date(year, month-1, 1, 0, 0, 0, 0, loc)
	startOfAfterNextMonth := time.Date(year, month+2, 1, 0, 0, 0, 0, loc)
	endOfNextMonth := startOfAfterNextMonth.AddDate(0, 0, -1)

	successCount := 0
	errorCount := 0
	totalDays := 0
	errorDetails := make(map[string]int)

	fmt.Printf("\n→ Mulai seeding data dari %s sampai %s\n",
		startOfPrevMonth.Format("2006-01-02"),
		endOfNextMonth.Format("2006-01-02"))

	// Seed random dengan waktu saat ini
	rand.Seed(time.Now().UnixNano())

	// 4. Loop tanggal & kirim data ke API
	for d := startOfPrevMonth; !d.After(endOfNextMonth); d = d.AddDate(0, 0, 1) {
		totalDays++

		// --- Nilai random untuk setiap hari ---
		basahLatex := float64(rand.Intn(2600-2000+1) + 2000)
		sheet := float64(rand.Intn(1200-900+1) + 900)
		basahLump := float64(rand.Intn(1200-900+1) + 900)
		brcr := float64(rand.Intn(1000-600+1) + 600)
		// --------------------------------------

		penyadap := map[string]interface{}{
			"IdBakuMandor": 43,
			"IdPenyadap":   1,
			"Tanggal":      d.Format("2006-01-02T15:04:05Z07:00"),
			"Tipe":         "BAKU_INTERNAL",
			"TahunTanam":   2020,
			"BasahLatex":   basahLatex,
			"Sheet":        sheet,
			"BasahLump":    basahLump,
			"BrCr":         brcr,
		}

		body, _ := json.Marshal(penyadap)
		req, _ := http.NewRequest("POST", "http://localhost:8080/api/baku", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{Name: "session_token", Value: SessionToken})

		res, err := client.Do(req)
		if err != nil {
			errorMsg := fmt.Sprintf("Network error: %v", err)
			fmt.Printf("  ✗ Tanggal %s: %s\n", d.Format("2006-01-02"), errorMsg)
			errorDetails[errorMsg]++
			errorCount++
			time.Sleep(500 * time.Millisecond)
			continue
		}

		if res.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(res.Body)
			errorMsg := fmt.Sprintf("Status %d: %s", res.StatusCode, string(bodyBytes))
			if errorDetails[errorMsg] == 0 {
				fmt.Printf("  ✗ Tanggal %s: %s\n", d.Format("2006-01-02"), errorMsg)
			}
			errorDetails[errorMsg]++
			errorCount++
		} else {
			successCount++
			if totalDays%5 == 0 || totalDays == 1 {
				fmt.Printf("  ✓ Progress: %d/%d hari berhasil\n", successCount, int(endOfNextMonth.Sub(startOfPrevMonth).Hours()/24)+1)
			}
		}
		res.Body.Close()

		time.Sleep(50 * time.Millisecond)
	}

	// 5. Summary
	fmt.Printf("\n→ Seeding data selesai: %d berhasil, %d gagal dari %d hari\n",
		successCount, errorCount, totalDays)

	if errorCount > 0 {
		fmt.Println("\n⚠️  DETAIL ERROR:")
		for errorMsg, count := range errorDetails {
			fmt.Printf("   - %s (terjadi %d kali)\n", errorMsg, count)
		}

		fmt.Println("\n💡 KEMUNGKINAN PENYEBAB & SOLUSI:")
		fmt.Println("   1. Mandor ID 1 tidak ada → Jalankan seed.SeedMandor()")
		fmt.Println("   2. Penyadap ID 1 tidak ada → Jalankan seed.SeedPenyadap()")
		fmt.Println("   3. Tipe 'BAKU' tidak valid → Cek models/baku.go")
		fmt.Println("   4. Session expired → Coba perpanjang waktu tunggu di main.go")
		fmt.Println("   5. Duplikasi data → Cek constraint unique di tabel baku")
	}
}
