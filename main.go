package main

import (
	"app-inputan-ptpn/config"
	"app-inputan-ptpn/controllers"
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

	// Routes
	setupRoutes()

	// Start server
	port := ":8080"
	fmt.Printf("Server started on http://localhost%s\n", port)
	fmt.Println("Login dengan: username=admin, password=admin123")

	log.Fatal(http.ListenAndServe(port, nil))
}

func setupRoutes() {
	// Static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// Auth routes
	http.HandleFunc("/", controllers.ServeLoginPage)
	http.HandleFunc("/login", controllers.Login)
	http.HandleFunc("/logout", controllers.Logout)

	// Protected routes (dengan middleware auth)
	http.HandleFunc("/dashboard", controllers.AuthMiddleware(serveDashboard))
	http.HandleFunc("/api/baku", controllers.AuthMiddleware(controllers.GetAllBaku))

	// API routes yang memerlukan autentikasi
	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		controllers.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"success": false, "message": "API endpoint not found"}`))
		})(w, r)
	})
}
