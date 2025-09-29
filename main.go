package main

import (
	"app-inputan-ptpn/config"
	"app-inputan-ptpn/routes"
	"app-inputan-ptpn/seed"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	// Banner aplikasi
	printBanner()

	// Hapus database lama jika ada
	dbFile := "produksi.db"
	if _, err := os.Stat(dbFile); err == nil {
		err := os.Remove(dbFile)
		if err != nil {
			log.Fatalf("❌ Gagal menghapus database lama: %v", err)
		}
		fmt.Println("✓ Database lama berhasil dihapus")
	}

	// Initialize database
	fmt.Println("\n🔧 Inisialisasi database...")
	config.InitDB()
	fmt.Println("✓ Database berhasil diinisialisasi")

	// Create templates directory if not exists
	if _, err := os.Stat("templates"); os.IsNotExist(err) {
		os.Mkdir("templates", 0755)
		fmt.Println("✓ Direktori templates dibuat")
	}

	// Create static directory if not exists
	if _, err := os.Stat("static"); os.IsNotExist(err) {
		os.Mkdir("static", 0755)
		fmt.Println("✓ Direktori static dibuat")
	}

	// Seed master data SEBELUM server jalan
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("📦 SEEDING MASTER DATA")
	fmt.Println(strings.Repeat("=", 60))

	fmt.Println("\n1️⃣  Seeding Mandor...")
	seed.SeedMandor()
	fmt.Println("   ✓ Mandor seeding selesai")

	fmt.Println("\n2️⃣  Seeding Penyadap...")
	seed.SeedPenyadap()
	fmt.Println("   ✓ Penyadap seeding selesai")

	fmt.Println("\n3️⃣  Seeding Baku (Data Awal)...")
	seed.SeedBaku()

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("✅ MASTER DATA SEEDING SELESAI")
	fmt.Println(strings.Repeat("=", 60))

	// Setup routes
	fmt.Println("\n🔧 Setup routing...")
	routes.SetupRoutes()
	fmt.Println("✓ Routing berhasil dikonfigurasi")

	// Start server di goroutine
	serverReady := make(chan bool)
	go func() {
		port := ":8080"
		fmt.Println("\n" + strings.Repeat("=", 60))
		fmt.Println("🚀 SERVER STARTING")
		fmt.Println(strings.Repeat("=", 60))
		fmt.Printf("   URL: http://localhost%s\n", port)
		fmt.Println("   Login credentials:")
		fmt.Println("   - Username: admin")
		fmt.Println("   - Password: admin123")
		fmt.Println(strings.Repeat("=", 60))

		// Signal bahwa server siap
		time.Sleep(1 * time.Second)
		serverReady <- true

		// Start server
		if err := http.ListenAndServe(port, nil); err != nil {
			log.Fatal("❌ Server error:", err)
		}
	}()

	// Tunggu server benar-benar siap
	<-serverReady
	fmt.Println("\n⏳ Menunggu server siap menerima request...")
	time.Sleep(5 * time.Second)
	fmt.Println("✓ Server siap!")

	// Validasi data master sebelum seeding transaksi
	if !seed.ValidateSeederData() {
		log.Fatal("❌ Validasi gagal. Seeding dibatalkan.")
	}

	// Jalankan seeder yang butuh API call
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("📊 SEEDING DATA TRANSAKSI (via API)")
	fmt.Println(strings.Repeat("=", 60))

	fmt.Println("\n1️⃣  Seeding Data Harian (1 Bulan)...")
	seed.SeedData()

	fmt.Println("\n2️⃣  Seeding Baku Borong...")
	seed.SeedBakuBorong()

	// Debug & Validasi hasil
	seed.DebugBakuPenyadapData()
	seed.CheckDuplicateEntries()

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("✅ SEEDING DATA TRANSAKSI SELESAI")
	fmt.Println(strings.Repeat("=", 60))

	// Tampilkan summary
	printSummary()

	// Keep main goroutine alive
	fmt.Println("\n✨ Aplikasi siap digunakan!")
	fmt.Println("   Tekan Ctrl+C untuk menghentikan server")
	fmt.Println()

	select {} // Block forever
}

// printBanner menampilkan banner aplikasi
func printBanner() {
	banner := `
╔══════════════════════════════════════════════════════════════╗
║                                                              ║
║     █████╗ ██████╗ ██████╗     ██████╗ ████████╗██████╗ ███╗   ██╗
║    ██╔══██╗██╔══██╗██╔══██╗    ██╔══██╗╚══██╔══╝██╔══██╗████╗  ██║
║    ███████║██████╔╝██████╔╝    ██████╔╝   ██║   ██████╔╝██╔██╗ ██║
║    ██╔══██║██╔═══╝ ██╔═══╝     ██╔═══╝    ██║   ██╔═══╝ ██║╚██╗██║
║    ██║  ██║██║     ██║         ██║        ██║   ██║     ██║ ╚████║
║    ╚═╝  ╚═╝╚═╝     ╚═╝         ╚═╝        ╚═╝   ╚═╝     ╚═╝  ╚═══╝
║                                                              ║
║              Sistem Input Data Produksi PTPN                 ║
║                      Version 1.0.0                           ║
║                                                              ║
╚══════════════════════════════════════════════════════════════╝
`
	fmt.Println(banner)
}

// printSummary menampilkan ringkasan data setelah seeding
func printSummary() {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("📈 RINGKASAN DATA")
	fmt.Println(strings.Repeat("=", 60))

	// Count records from database
	var mandorCount, penyadapCount, bakuPenyadapCount, bakuDetailCount int64

	config.DB.Model(&struct {
		ID uint `gorm:"primaryKey"`
	}{}).Table("baku_mandors").Count(&mandorCount)

	config.DB.Model(&struct {
		ID uint `gorm:"primaryKey"`
	}{}).Table("penyadaps").Count(&penyadapCount)

	config.DB.Model(&struct {
		ID uint `gorm:"primaryKey"`
	}{}).Table("baku_penyadaps").Count(&bakuPenyadapCount)

	config.DB.Model(&struct {
		ID uint `gorm:"primaryKey"`
	}{}).Table("baku_details").Count(&bakuDetailCount)

	fmt.Printf("   📋 Total Mandor        : %d records\n", mandorCount)
	fmt.Printf("   👥 Total Penyadap      : %d records\n", penyadapCount)
	fmt.Printf("   📊 Total Baku Penyadap : %d records\n", bakuPenyadapCount)
	fmt.Printf("   📑 Total Baku Detail   : %d records\n", bakuDetailCount)
	fmt.Println(strings.Repeat("=", 60))
}

// Import strings untuk strings.Repeat
