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

	// ================== BAKU PAGE (HTML) ==================
	protected.HandleFunc("/baku", controllers.ServeBakuPage).Methods("GET")

	// ================== NEW: TIPE PRODUKSI API ==================
	protected.HandleFunc("/api/tipe-produksi", controllers.GetTipeProduksiList).Methods("GET")

	// ================== BAKU API (CRUD JSON) ==================
	// PENTING: Route spesifik HARUS didefinisikan SEBELUM route dengan parameter!

	// BAKU DETAIL API - HARUS SEBELUM /api/baku/{id}
	// Support query parameter: ?tipe=BAKU|BAKU_BORONG|BORONG_EXTERNAL|BORONG_INTERNAL|TETES_LANJUT
	protected.HandleFunc("/api/baku/detail", controllers.GetAllBakuDetail).Methods("GET")
	protected.HandleFunc("/api/baku/detail/{tanggal}", controllers.GetBakuDetailByDate).Methods("GET")

	// ================== BAKU PENYADAP CRUD - Setelah route detail ==================
	// Support query parameter: ?tipe=BAKU|BAKU_BORONG|BORONG_EXTERNAL|BORONG_INTERNAL|TETES_LANJUT
	protected.HandleFunc("/api/baku", controllers.GetAllBakuPenyadap).Methods("GET")
	protected.HandleFunc("/api/baku", controllers.CreateBakuPenyadap).Methods("POST")
	protected.HandleFunc("/api/baku/{id}", controllers.GetBakuPenyadapByID).Methods("GET") // pastikan ada di controller
	protected.HandleFunc("/api/baku/{id}", controllers.UpdateBakuPenyadap).Methods("PUT")
	protected.HandleFunc("/api/baku/{id}", controllers.DeleteBakuPenyadap).Methods("DELETE") // ganti Delete -> DeleteBakuPenyadap

	// ================== MANDOR API (CRUD) ==================
	protected.HandleFunc("/api/mandor", controllers.GetAllMandor).Methods("GET")
	protected.HandleFunc("/api/mandor", controllers.CreateMandor).Methods("POST")
	protected.HandleFunc("/api/mandor/{id}", controllers.GetMandorByID).Methods("GET")
	protected.HandleFunc("/api/mandor/{id}", controllers.UpdateMandor).Methods("PUT")
	protected.HandleFunc("/api/mandor/{id}", controllers.DeleteMandor).Methods("DELETE")

	// ================== REPORTING API WITH TIPE SUPPORT ==================
	// Summary Mandor (total dari semua penyadap)
	// Support query parameter: ?tipe=BAKU|BAKU_BORONG|BORONG_EXTERNAL|BORONG_INTERNAL|TETES_LANJUT
	protected.HandleFunc("/api/reporting/mandor", controllers.GetMandorSummaryAll).Methods("GET")
	protected.HandleFunc("/api/reporting/mandor/{tanggal}", controllers.GetMandorSummaryByDate).Methods("GET")

	// Detail individual penyadap
	// Support query parameter: ?tipe=BAKU|BAKU_BORONG|BORONG_EXTERNAL|BORONG_INTERNAL|TETES_LANJUT
	protected.HandleFunc("/api/reporting/penyadap", controllers.GetPenyadapDetailAll).Methods("GET")
	protected.HandleFunc("/api/reporting/penyadap/{tanggal}", controllers.GetPenyadapDetailByDate).Methods("GET")

	// ================== SEARCH API WITH TIPE SUPPORT ==================
	// Pencarian mandor berdasarkan nama dengan filter tanggal dan tipe opsional
	// Query parameters: ?nama=xxx&tanggal=YYYY-MM-DD&tipe=BAKU
	protected.HandleFunc("/api/search/mandor", controllers.SearchMandorByName).Methods("GET")

	// Pencarian penyadap berdasarkan nama dengan filter tanggal dan tipe opsional
	// Query parameters: ?nama=xxx&tanggal=YYYY-MM-DD&tipe=BAKU
	protected.HandleFunc("/api/search/penyadap", controllers.SearchPenyadapByName).Methods("GET")

	// Detail mandor beserta semua penyadapnya dengan filter tipe
	// Query parameters: ?nama=xxx&tanggal=YYYY-MM-DD&tipe=BAKU
	protected.HandleFunc("/api/search/mandor/detail", controllers.GetMandorWithPenyadapDetail).Methods("GET")

	// Global search dengan berbagai filter termasuk tipe
	// Query parameters: ?type=mandor|penyadap&nama=xxx&tanggal=YYYY-MM-DD&afdeling=xxx&tahun=2024&tipe=BAKU
	protected.HandleFunc("/api/search/all", controllers.SearchAll).Methods("GET")

	// ================== PENYADAP API (CRUD + Search) ==================
	// Search HARUS sebelum route dengan parameter {id}
	protected.HandleFunc("/api/penyadap/search", controllers.GetPenyadapByName).Methods("GET")
	protected.HandleFunc("/api/penyadap", controllers.GetAllPenyadap).Methods("GET")
	protected.HandleFunc("/api/penyadap", controllers.CreatePenyadap).Methods("POST")
	protected.HandleFunc("/api/penyadap/{id}", controllers.UpdatePenyadap).Methods("PUT")
	protected.HandleFunc("/api/penyadap/{id}", controllers.DeletePenyadap).Methods("DELETE")

	// ================== BACKWARD COMPATIBILITY ==================
	// Keep old routes if needed
	protected.HandleFunc("/mandor", controllers.GetAllMandor).Methods("GET")
	protected.HandleFunc("/mandor", controllers.CreateMandor).Methods("POST")
	protected.HandleFunc("/mandor/{id}", controllers.UpdateMandor).Methods("PUT")
	protected.HandleFunc("/mandor/{id}", controllers.DeleteMandor).Methods("DELETE")
	protected.HandleFunc("/penyadap/search", controllers.GetPenyadapByName).Methods("GET")

	// ================== ERROR HANDLING ==================
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
