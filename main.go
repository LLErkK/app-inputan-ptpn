package main

import (
	"app-inputan-ptpn/config"
	"app-inputan-ptpn/routes"
	"app-inputan-ptpn/seed"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
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

	// Start server
	port := ":8080"
	fmt.Printf("Server started on http://localhost%s\n", port)
	fmt.Println("Login dengan: username=admin, password=admin123")

	log.Fatal(http.ListenAndServe(port, nil))
}
