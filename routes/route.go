package routes

import (
	"app-inputan-ptpn/controllers"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"
)

// SetupRoutes mengatur semua routing aplikasi dengan dukungan filter range tanggal
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

	// Serve KML files
	r.PathPrefix("/kml/").Handler(http.StripPrefix("/kml/",
		http.FileServer(http.Dir("./templates/kml/"))))

	// ================== AUTH ROUTES ==================
	r.HandleFunc("/", controllers.ServeLoginPage).Methods("GET")
	r.HandleFunc("/login", controllers.Login).Methods("POST")
	r.HandleFunc("/logout", controllers.Logout).Methods("GET")

	// ================== PROTECTED ROUTES ==================
	protected := r.PathPrefix("/").Subrouter()
	protected.Use(authMiddleware)

	// Dashboard
	protected.HandleFunc("/dashboard", controllers.ServeDashboardPage).Methods("GET")
	protected.HandleFunc("/api/dashboard", controllers.GetDashboardData).Methods("GET")

	// ================== BAKU PAGE (HTML) ==================
	protected.HandleFunc("/baku", controllers.ServeBakuPage).Methods("GET")

	// ================== NEW: TIPE PRODUKSI API ==================
	protected.HandleFunc("/api/tipe-produksi", controllers.GetTipeProduksiList).Methods("GET")

	// ================== ENHANCED BAKU DETAIL API WITH DATE RANGE ==================
	protected.HandleFunc("/api/baku/detail", controllers.GetAllBakuDetail).Methods("GET")
	protected.HandleFunc("/api/baku/detail/range", controllers.GetBakuDetailByDateRange).Methods("GET")
	protected.HandleFunc("/api/baku/detail/{tanggal}", controllers.GetBakuDetailByDate).Methods("GET")

	// ================== ENHANCED BAKU PENYADAP CRUD WITH DATE RANGE ==================
	protected.HandleFunc("/api/baku", controllers.GetAllBakuPenyadap).Methods("GET")
	protected.HandleFunc("/api/baku", controllers.CreateBakuPenyadap).Methods("POST")
	protected.HandleFunc("/api/baku/{id}", controllers.GetBakuPenyadapByID).Methods("GET")
	protected.HandleFunc("/api/baku/{id}", controllers.UpdateBakuPenyadap).Methods("PUT")
	protected.HandleFunc("/api/baku/{id}", controllers.DeleteBakuPenyadap).Methods("DELETE")
	protected.HandleFunc("/api/baku/rekap/today", controllers.GetBakuPenyadapToday).Methods("GET")

	// ================== MANDOR API (CRUD) ==================
	protected.HandleFunc("/api/mandor/search", controllers.GetMandorByName).Methods("GET")
	protected.HandleFunc("/api/mandor", controllers.GetAllMandor).Methods("GET")
	protected.HandleFunc("/api/mandor", controllers.CreateMandor).Methods("POST")
	protected.HandleFunc("/api/mandor/{id}", controllers.GetMandorByID).Methods("GET")
	protected.HandleFunc("/api/mandor/{id}", controllers.UpdateMandor).Methods("PUT")
	protected.HandleFunc("/api/mandor/{id}", controllers.DeleteMandor).Methods("DELETE")

	// ================== ENHANCED REPORTING API WITH DATE RANGE SUPPORT ==================
	protected.HandleFunc("/api/reporting/mandor", controllers.GetMandorSummaryAll).Methods("GET")
	protected.HandleFunc("/api/reporting/mandor/range", controllers.GetMandorSummaryByDateRange).Methods("GET")
	protected.HandleFunc("/api/reporting/mandor/{tanggal}", controllers.GetMandorSummaryByDate).Methods("GET")

	protected.HandleFunc("/api/reporting/penyadap", controllers.GetPenyadapDetailAll).Methods("GET")
	protected.HandleFunc("/api/reporting/penyadap/range", controllers.GetPenyadapDetailByDateRange).Methods("GET")
	protected.HandleFunc("/api/reporting/penyadap/{tanggal}", controllers.GetPenyadapDetailByDate).Methods("GET")

	// ================== ENHANCED SEARCH API WITH DATE RANGE SUPPORT ==================
	protected.HandleFunc("/api/search", controllers.SearchData).Methods("GET")
	//parameternya idMandor/idPenyadap, tanggalAwal, tanggalAkhir, tipeProduksi

	// ================== PENYADAP API (CRUD + Search) ==================
	// Search HARUS sebelum route dengan parameter {id}
	protected.HandleFunc("/api/penyadap/search", controllers.GetPenyadapByName).Methods("GET")
	protected.HandleFunc("/api/penyadap", controllers.GetAllPenyadap).Methods("GET")
	protected.HandleFunc("/api/penyadap", controllers.CreatePenyadap).Methods("POST")
	protected.HandleFunc("/api/penyadap/{id}", controllers.UpdatePenyadap).Methods("PUT")
	protected.HandleFunc("/api/penyadap/{id}", controllers.DeletePenyadap).Methods("DELETE")

	//monitoring
	protected.HandleFunc("/monitoring", controllers.ServeMonitoringPage).Methods("GET")

	// ================== SMART MONITORING SEARCH (NEW ROUTE) ==================
	// Endpoint yang dipakai oleh halaman monitoring untuk pencarian pintar
	protected.HandleFunc("/api/monitoring/smart-search", controllers.SmartMonitoringSearch).Methods("GET")

	// ================== BACKWARD COMPATIBILITY ==================
	protected.HandleFunc("/mandor", controllers.GetAllMandor).Methods("GET")
	protected.HandleFunc("/mandor", controllers.CreateMandor).Methods("POST")
	protected.HandleFunc("/mandor/{id}", controllers.UpdateMandor).Methods("PUT")
	protected.HandleFunc("/mandor/{id}", controllers.DeleteMandor).Methods("DELETE")
	protected.HandleFunc("/penyadap/search", controllers.GetPenyadapByName).Methods("GET")

	// Routes untuk visualisasi
	protected.HandleFunc("/visualisasi", controllers.ServeVisualisasiPage).Methods("GET")
	protected.HandleFunc("/api/visualisasi", controllers.GetVisualisasiData).Methods("GET")

	//rekap endpoint
	protected.HandleFunc("/rekap", controllers.ServeRekapPage).Methods("GET")
	protected.HandleFunc("/rekap/today", controllers.GetBakuDetailToday).Methods("GET")
	protected.HandleFunc("/rekap/until-today", controllers.GetBakuDetailUntilTodayThisMonth).Methods("GET")

	//upload excell
	protected.HandleFunc("/upload", controllers.ServeUploadPage).Methods("GET")
	protected.HandleFunc("/api/upload", controllers.CreateUpload).Methods("POST")
	protected.HandleFunc("/api/upload", controllers.GetAllUploads).Methods("GET")
	protected.HandleFunc("/api/upload/range", controllers.GetUploadsByDateRange).Methods("GET")
	protected.HandleFunc("/api/upload/{id}", controllers.GetUploadByID).Methods("GET")
	protected.HandleFunc("/api/upload/{id}", controllers.DeleteUpload).Methods("DELETE")
	protected.HandleFunc("/api/upload/{id}/download", controllers.DownloadFile).Methods("GET")

	protected.HandleFunc("/api/master", controllers.GetAllMaster).Methods("GET")
	protected.HandleFunc("/api/master/{masterId}", controllers.DeleteMaster).Methods("DELETE") //

	// ================== ADDITIONAL MONITORING ENDPOINTS ==================
	protected.HandleFunc("/api/monitoring/today/summary", func(w http.ResponseWriter, r *http.Request) {
		// Redirect ke smart search dengan tanggal hari ini
		tipe := r.URL.Query().Get("tipe")
		today := time.Now().Format("2006-01-02")

		params := url.Values{}
		params.Add("tanggalAwal", today)
		params.Add("tanggalAkhir", today)
		if tipe != "" {
			params.Add("filterJenis", tipe)
		}

		http.Redirect(w, r, "/api/monitoring/smart-search?"+params.Encode(), http.StatusTemporaryRedirect)
	}).Methods("GET")

	protected.HandleFunc("/api/monitoring/week/summary", func(w http.ResponseWriter, r *http.Request) {
		tipe := r.URL.Query().Get("tipe")
		today := time.Now()

		// Cari hari Senin minggu ini
		offset := int(today.Weekday()) - int(time.Monday)
		if offset < 0 {
			offset = 6 // kalau hari ini Minggu (0), mundur ke Senin minggu sebelumnya
		}
		weekStart := today.AddDate(0, 0, -offset)

		// Akhir minggu = Minggu
		weekEnd := weekStart.AddDate(0, 0, 6)

		params := url.Values{}
		params.Add("tanggalAwal", weekStart.Format("2006-01-02"))
		params.Add("tanggalAkhir", weekEnd.Format("2006-01-02"))
		if tipe != "" {
			params.Add("filterJenis", tipe)
		}

		http.Redirect(w, r, "/api/monitoring/smart-search?"+params.Encode(), http.StatusTemporaryRedirect)
	}).Methods("GET")

	protected.HandleFunc("/api/monitoring/month/summary", func(w http.ResponseWriter, r *http.Request) {
		tipe := r.URL.Query().Get("tipe")
		today := time.Now()
		monthStart := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())

		params := url.Values{}
		params.Add("tanggalAwal", monthStart.Format("2006-01-02"))
		params.Add("tanggalAkhir", today.Format("2006-01-02"))
		if tipe != "" {
			params.Add("filterJenis", tipe)
		}

		http.Redirect(w, r, "/api/monitoring/smart-search?"+params.Encode(), http.StatusTemporaryRedirect)
	}).Methods("GET")

	// ================== ERROR HANDLING ==================
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
