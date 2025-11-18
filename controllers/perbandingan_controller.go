package controllers

import (
	"net/http"
)

func ServePerbandinganPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/html/perbandingan.html")
}
