package middleware

import (
	"context-monster/models"
	"net/http"
	"os"
	"strings"
)

var prefix = os.Getenv("path")

var noAuthPaths = []string{prefix + "apikeys/new"}

func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, value := range noAuthPaths {
			if value == r.URL.Path {
				next.ServeHTTP(w, r)
				return
			}
		}
		tokenHeader := r.Header.Get("Authorization")
		if tokenHeader == "" {
			unauthorized(w, "Missing Auth Token! Use Bearer Token to authenticate", http.StatusUnauthorized)
			return
		}
		split := strings.Split(tokenHeader, " ")
		if len(split) != 2 {
			unauthorized(w, "Missing Auth Token! Use Bearer Token to authenticate", http.StatusUnauthorized)
			return
		}
		token := split[1]
		apiKey, err := models.GetAPIKey(token)
		if err != nil {
			unauthorized(w, "Auth Token is not valid!", http.StatusUnauthorized)
			return
		}
		if apiKey.AccessLevel == models.AccessNothing ||
			r.Method == "GET" && !apiKey.CanRead() || r.Method == "POST" && !apiKey.CanWrite() {
			apiKey.RegisterUse(r.URL.Path, r.Method, false)
			unauthorized(w, "You are unauthorized to perform this action.", http.StatusForbidden)
			return
		}
		apiKey.RegisterUse(r.URL.Path, r.Method, true)
		next.ServeHTTP(w, r)
	})
}

func unauthorized(w http.ResponseWriter, msg string, statusCode int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write([]byte(`{"message": "` + msg + `"}`))
}
