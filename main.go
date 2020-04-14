package main

import (
	"context-monster/controllers"
	"context-monster/middleware"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

var prefix = os.Getenv("path")

func main() {

	router := mux.NewRouter()

	router.HandleFunc(prefix+"apikeys/new", controllers.CreateAPIKey).Methods("GET")
	router.HandleFunc(prefix+"apikeys/info", controllers.ViewAPIKey).Methods("GET")
	router.HandleFunc(prefix+"apikeys/logs", controllers.ViewAPIKeyLogs).Methods("GET")

	router.HandleFunc(prefix+"tweets/{id}", controllers.GetTweet).Methods("GET")
	router.HandleFunc(prefix+"tweets", controllers.GetTweets).Methods("GET")
	router.HandleFunc(prefix+"tweets/search/{needle}", controllers.SearchTweets).Methods("GET")
	router.HandleFunc(prefix+"tweets", controllers.CreateTweet).Methods("POST")

	router.HandleFunc(prefix+"keywords/{text}", controllers.GetKeyword).Methods("GET")
	router.HandleFunc(prefix+"keywords", controllers.GetKeywords).Methods("GET")
	router.HandleFunc(prefix+"keywords/search/{needle}", controllers.SearchKeywords).Methods("GET")

	router.NotFoundHandler = http.HandlerFunc(controllers.NotFound404Error)

	router.Use(middleware.Authentication) // Auth middleware

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" //localhost
	}

	fmt.Printf("Running on \033[1;35m%s\n\033[0m", port)

	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		fmt.Print(err)
	}
}
