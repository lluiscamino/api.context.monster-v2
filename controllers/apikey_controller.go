package controllers

import (
	"context-monster/models"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"
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

func ViewAPIKey(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	token, err := getToken(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "Unexpected error"}`))
		return
	}
	apiKey, err := models.GetAPIKey(token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "Unexpected error"}`))
		return
	}
	json.NewEncoder(w).Encode(apiKey)
}

func ViewAPIKeyLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	token, err := getToken(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "Unexpected error"}`))
		return
	}
	apiKey, err := models.GetAPIKey(token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "Unexpected error"}`))
		return
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit > maxLimit {
		limit = maxLimit
	}
	if limit < 0 {
		limit = 0
	}
	apiKey.LoadLogs(uint(limit))
	json.NewEncoder(w).Encode(apiKey.Logs)
}

func getToken(r *http.Request) (string, error) {
	tokenHeader := r.Header.Get("Authorization")
	if tokenHeader == "" {
		return "", errors.New("missing Authorization header")
	}
	split := strings.Split(tokenHeader, " ")
	if len(split) != 2 {
		return "", errors.New("missing Authorization token")
	}
	return split[1], nil
}
