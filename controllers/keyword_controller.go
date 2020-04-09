package controllers

import (
	"context-monster/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func GetKeyword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	text := mux.Vars(r)["text"]
	keyword := models.GetKeyword(text)
	if keyword == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "Keyword ` + text + ` does not exist"}`))
		return
	}
	json.NewEncoder(w).Encode(keyword)
}

func GetKeywords(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 100
	}
	loadRatings, err := strconv.ParseBool(r.URL.Query().Get("ratings"))
	if err != nil {
		loadRatings = true
	}
	orderBy := r.URL.Query().Get("order")
	if orderBy != "counter" && orderBy != "searches" {
		orderBy = "counter"
	}
	keywords := models.GetKeywords(uint(limit), orderBy, loadRatings)
	json.NewEncoder(w).Encode(keywords)
}

func SearchKeywords(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 100
	}
	loadRatings, err := strconv.ParseBool(r.URL.Query().Get("ratings"))
	if err != nil {
		loadRatings = true
	}
	needle := mux.Vars(r)["needle"]
	keywords := models.SearchKeywords(needle, uint(limit), loadRatings)
	json.NewEncoder(w).Encode(keywords)
}
