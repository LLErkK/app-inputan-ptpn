package routes

import (
	"app-inputan-ptpn/controllers"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// SetupRoutes mengatur semua routing aplikasi
func SetupRoutes() {
	r := mux.NewRouter()

	// PENTING: Static files harus SEBELUM protected routes
	// Agar tidak terblokir oleh auth middleware
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir("./templates/"))))

	// Debug: check templates directory
	if _, err := os.Stat("./templates"); os.IsNotExist(err) {
		log.Fatal("Templates directory does not exist!")
	}

	// Auth routes (tidak perlu auth)
	r.HandleFunc("/", controllers.ServeLoginPage).Methods("GET")
	r.HandleFunc("/login", controllers.Login).Methods("POST")
	r.HandleFunc("/logout", controllers.Logout).Methods("GET")

	// Protected routes
	protected := r.PathPrefix("/").Subrouter()
	protected.Use(authMiddleware)

	// Dashboard
	protected.HandleFunc("/dashboard", controllers.ServeDashboardPage).Methods("GET")

	// Baku page
	protected.HandleFunc("/baku", controllers.ServeBakuPage).Methods("GET")

	// API Routes - Mandor CRUD
	protected.HandleFunc("/mandor", controllers.GetAllMandor).Methods("GET")
	protected.HandleFunc("/mandor", controllers.CreateMandor).Methods("POST")
	protected.HandleFunc("/mandor/{id}", controllers.UpdateMandor).Methods("PUT")
	protected.HandleFunc("/mandor/{id}", controllers.DeleteMandor).Methods("DELETE")

	// API Routes - Penyadap CRUD
	protected.HandleFunc("/penyadap", controllers.GetAllBakuPenyadap).Methods("GET")
	protected.HandleFunc("/penyadap", controllers.CreateBakuPenyadap).Methods("POST")
	protected.HandleFunc("/penyadap/{id}", controllers.GetBakuPenyadapByID).Methods("GET")
	protected.HandleFunc("/penyadap/{id}", controllers.UpdateBakuPenyadap).Methods("PUT")
	protected.HandleFunc("/penyadap/{id}", controllers.DeleteBakuPenyadap).Methods("DELETE")

	http.Handle("/", r)
}

// Middleware wrapper
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.AuthMiddleware(next.ServeHTTP)(w, r)
	})
}
