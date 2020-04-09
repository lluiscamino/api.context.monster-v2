package controllers

import (
	"context-monster/models"
	"net/http"
	"os"
)

func CreateAPIKey(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if os.Getenv("enable_apikey_creation") != "true" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"message": "This feature is currently disabled on the server."}`))
		return
	}
	name := r.URL.Query().Get("name")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "Please provide a APIKey name."}`))
		return
	}
	token, err := models.CreateAPIKey(name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "API Key could not due to an unexpected error."}`))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"token": "` + token + `"}`))
}
