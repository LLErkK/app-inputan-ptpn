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

	http.HandleFunc("/baku", controllers.AuthMiddleware(controllers.ServeBakuPage))
	// Catch-all untuk API yang tidak ditemukan
	http.HandleFunc("/api/", controllers.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"success": false, "message": "API endpoint not found"}`))
	}))
}
