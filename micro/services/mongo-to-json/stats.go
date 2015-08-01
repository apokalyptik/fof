package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func allTimeStats(w http.ResponseWriter, r *http.Request) {
	setCORS(w, r)
	setJSON(w, r)
	vars := mux.Vars(r)
	if data, err := atsCache.get(vars["section"], vars["stat"]); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Write(data)
	}
}

func allTimeStatKeys(w http.ResponseWriter, r *http.Request) {
	setCORS(w, r)
	setJSON(w, r)
	if data, err := atkCache.get(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Write(data)
	}
}
