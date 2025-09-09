package routes

import (
	"app-inputan-ptpn/controllers"
	"github.com/gorilla/mux"
	"net/http"
)

// SetupRoutes mengatur semua routing aplikasi
func SetupRoutes() {
	r := mux.NewRouter()

	// ================== STATIC FILES ==================
	// Serve CSS
	r.PathPrefix("/css/").Handler(http.StripPrefix("/css/",
		http.FileServer(http.Dir("./templates/css/"))))

	// Serve JS
	r.PathPrefix("/js/").Handler(http.StripPrefix("/js/",
		http.FileServer(http.Dir("./templates/js/"))))

	// Serve Assets (gambar, ikon, dll)
	r.PathPrefix("/asset/").Handler(http.StripPrefix("/asset/",
		http.FileServer(http.Dir("./templates/asset/"))))

	// ================== AUTH ROUTES ==================
	r.HandleFunc("/", controllers.ServeLoginPage).Methods("GET")
	r.HandleFunc("/login", controllers.Login).Methods("POST")
	r.HandleFunc("/logout", controllers.Logout).Methods("GET")

	// ================== PROTECTED ROUTES ==================
	protected := r.PathPrefix("/").Subrouter()
	protected.Use(authMiddleware)

	// Dashboard
	protected.HandleFunc("/dashboard", controllers.ServeDashboardPage).Methods("GET")

	// Baku
	// ================== BAKU PAGE (HTML) ==================
	protected.HandleFunc("/baku", controllers.ServeBakuPage).Methods("GET")

	// ================== BAKU API (CRUD JSON) ==================
	protected.HandleFunc("/api/baku", controllers.GetAllBakuPenyadap).Methods("GET")
	protected.HandleFunc("/api/baku/{id}", controllers.GetBakuPenyadapByID).Methods("GET")
	protected.HandleFunc("/api/baku", controllers.CreateBakuPenyadap).Methods("POST")
	protected.HandleFunc("/api/baku/{id}", controllers.UpdateBakuPenyadap).Methods("PUT")
	protected.HandleFunc("/api/baku/{id}", controllers.DeleteBakuPenyadap).Methods("DELETE")

	// Rekap / BakuDetail
	protected.HandleFunc("/api/baku/detail", controllers.GetAllBakuDetail).Methods("GET")
	protected.HandleFunc("/api/baku/detail/{tanggal}", controllers.GetBakuDetailByDate).Methods("GET")
	// Mandor CRUD
	protected.HandleFunc("/mandor", controllers.GetAllMandor).Methods("GET")
	protected.HandleFunc("/mandor", controllers.CreateMandor).Methods("POST")
	protected.HandleFunc("/mandor/{id}", controllers.UpdateMandor).Methods("PUT")
	protected.HandleFunc("/mandor/{id}", controllers.DeleteMandor).Methods("DELETE")

	// Penyadap CRUD + Search
	protected.HandleFunc("/api/penyadap", controllers.GetAllPenyadap).Methods("GET")
	protected.HandleFunc("/api/penyadap", controllers.CreatePenyadap).Methods("POST")
	protected.HandleFunc("/api/penyadap/{id}", controllers.UpdatePenyadap).Methods("PUT")
	protected.HandleFunc("/api/penyadap/{id}", controllers.DeletePenyadap).Methods("DELETE")
	protected.HandleFunc("/api/penyadap/search", controllers.GetPenyadapByName).Methods("GET")

	// Search Penyadap
	protected.HandleFunc("/penyadap/search", controllers.GetPenyadapByName).Methods("GET")

	// Catch-all untuk API yang tidak ditemukan
	protected.PathPrefix("/api/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"success": false, "message": "API endpoint not found"}`))
	})

	http.Handle("/", r)
}

// Middleware wrapper
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.AuthMiddleware(next.ServeHTTP)(w, r)
	})
}
