package main

import (
	"app-inputan-ptpn/config"
	"app-inputan-ptpn/routes"
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
	time.Sleep(2 * time.Second)
	fmt.Println("✓ Server siap!")

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
