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
	protected.HandleFunc("/baku", controllers.ServeBakuPage).Methods("GET")
	protected.HandleFunc("/baku", controllers.CreateBakuPenyadap).Methods("POST")

	// Mandor CRUD
	protected.HandleFunc("/mandor", controllers.GetAllMandor).Methods("GET")
	protected.HandleFunc("/mandor", controllers.CreateMandor).Methods("POST")
	protected.HandleFunc("/mandor/{id}", controllers.UpdateMandor).Methods("PUT")
	protected.HandleFunc("/mandor/{id}", controllers.DeleteMandor).Methods("DELETE")

	// Penyadap CRUD
	protected.HandleFunc("/penyadap", controllers.GetAllBakuPenyadap).Methods("GET")
	protected.HandleFunc("/penyadap", controllers.CreateBakuPenyadap).Methods("POST")
	protected.HandleFunc("/penyadap/{id}", controllers.GetBakuPenyadapByID).Methods("GET")
	protected.HandleFunc("/penyadap/{id}", controllers.UpdateBakuPenyadap).Methods("PUT")
	protected.HandleFunc("/penyadap/{id}", controllers.DeleteBakuPenyadap).Methods("DELETE")

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
