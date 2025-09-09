package routes

import (
	"app-inputan-ptpn/controllers"
	"net/http"
)

// SetupRoutes mengatur semua routing aplikasi
func SetupRoutes() {
	// Static files - pastikan path benar
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("templates/"))))

	// Auth routes (tidak butuh middleware)
	http.HandleFunc("/", controllers.ServeLoginPage)
	http.HandleFunc("/login", controllers.Login)
	http.HandleFunc("/logout", controllers.Logout)

	// Protected routes (butuh middleware)
	ProtectedRoutes()
}

func ProtectedRoutes() {
	// Dashboard
	http.HandleFunc("/dashboard", controllers.AuthMiddleware(controllers.ServeDashboardPage))
	//roote baku
	http.HandleFunc("/baku", controllers.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.ServeBakuPage(w, r)
		case http.MethodPost:
			controllers.CreateBakuPenyadap(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/mandor", controllers.AuthMiddleware(
		func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				controllers.GetAllMandor(w, r)
			case http.MethodPost:
				controllers.CreateMandor(w, r)
			case http.MethodDelete:
				controllers.DeleteMandor(w, r)
			case http.MethodPut:
				controllers.UpdateMandor(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}

		}))
	// Catch-all untuk API yang tidak ditemukan
	http.HandleFunc("/api/", controllers.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"success": false, "message": "API endpoint not found"}`))
	}))
}
