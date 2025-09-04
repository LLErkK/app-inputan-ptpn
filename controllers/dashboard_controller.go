package controllers

import "net/http"

func ServeDashboardPage(w http.ResponseWriter, r *http.Request) {

	// Serve login HTML file
	http.ServeFile(w, r, "templates/html/login.html")
}
