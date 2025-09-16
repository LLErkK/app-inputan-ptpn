package controllers

import (
	"html/template"
	"net/http"
)

func ServeDashboardPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/html/dashboard.html")
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}
