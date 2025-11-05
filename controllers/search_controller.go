package controllers

import "net/http"

func SearchData(w http.ResponseWriter, r *http.Request) {
	idMandor := r.URL.Query().Get("idMandor")
	idPenyadap := r.URL.Query().Get("idPenyadap")

	if idMandor == "" && idPenyadap == "" {
		http.Error(w, "tidak ada parameter id", http.StatusBadRequest)
		return
	}
	if idPenyadap != "" && idMandor == "" {
		SearchPenyadap(w, r)
	} else {
		SearchMandor(w, r)
	}

}
