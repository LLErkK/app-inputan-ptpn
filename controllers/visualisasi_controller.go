package controllers

import "net/http"

func ServeVisualisasiPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/html/visualisasi.html")
}
