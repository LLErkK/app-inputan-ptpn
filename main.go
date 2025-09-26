package main

import (
	"app-inputan-ptpn/config"
	"app-inputan-ptpn/routes"
	"app-inputan-ptpn/seed"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	dbFile := "produksi.db"
	if _, err := os.Stat(dbFile); err == nil {
		err := os.Remove(dbFile)
		if err != nil {
			log.Fatalf("Gagal menghapus database lama: %v", err)
		}
		fmt.Println("Database lama berhasil dihapus")
	}

	// Initialize database
	config.InitDB()

	// Create templates directory if not exists
	if _, err := os.Stat("templates"); os.IsNotExist(err) {
		os.Mkdir("templates", 0755)
	}

	// Create static directory if not exists
	if _, err := os.Stat("static"); os.IsNotExist(err) {
		os.Mkdir("static", 0755)
	}

	// Setup routes
	routes.SetupRoutes()
	seed.SeedMandor()
	seed.SeedPenyadap()
	seed.SeedBaku()

	// Jalankan seeder setelah server aktif
	go func() {
		time.Sleep(2 * time.Second)
		seed.SeedData()
		seed.SeedBakuBorong()
	}()

	// Start server
	port := ":8080"
	fmt.Printf("Server started on http://localhost%s\n", port)
	fmt.Println("Login dengan: username=admin, password=admin123")

	log.Fatal(http.ListenAndServe(port, nil))
}
