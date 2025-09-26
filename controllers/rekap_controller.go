package controllers

import "net/http"

func ServeRekapPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/html/rekap.html")
}
