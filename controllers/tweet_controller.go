package controllers

import (
	"context-monster/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func CreateTweet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tweet := &models.Tweet{}
	if err := json.NewDecoder(r.Body).Decode(tweet); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "Tweet could not be created due to an incorrect format"}`))
		return
	}
	if len(tweet.Ratings) > 0 {
		tweet.Title = tweet.Ratings[0].KeywordText
	} else {
		tweet.Title = "Undefined"
	}
	tweet.NumKeywords = uint(len(tweet.Ratings))
	tweet.Date = time.Now()
	tweet.ID = 0 // Ignore ID
	for i, r := range tweet.Ratings {
		tweet.Ratings[i].Keyword = &models.Keyword{
			Text: r.KeywordText,
		}
	}
	if err := tweet.Create(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "Tweet could not be created: ` + err.Error() + `"}`))
		return
	}
	for i := range tweet.Ratings {
		if err := tweet.Ratings[i].Keyword.CreateOrUpdate(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"message": "Invalid keyword: ` + err.Error() + `"}`))
			return
		}
		tweet.Ratings[i].TweetID = 0
		tweet.Ratings[i].KeywordText = ""
	}
	tweet.IsFirst = tweet.ID == 1
	tweet.IsLast = true
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tweet)
}

func GetTweet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "Please provide a valid Tweet ID"}`))
		return
	}
	tweet := models.GetTweet(uint(id))
	if tweet == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "Tweet with ID ` + strconv.Itoa(id) + ` does not exist"}`))
		return
	}
	json.NewEncoder(w).Encode(tweet)
}

func GetTweets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 100
	}
	loadRatings, err := strconv.ParseBool(r.URL.Query().Get("ratings"))
	if err != nil {
		loadRatings = true
	}
	order := strings.ToUpper(r.URL.Query().Get("order"))
	if order != "ASC" && order != "DESC" {
		order = "ASC"
	}
	tweets := models.GetTweets(uint(limit), order, loadRatings)
	json.NewEncoder(w).Encode(tweets)
}

func SearchTweets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 100
	}
	order := strings.ToUpper(r.URL.Query().Get("order"))
	if order != "ASC" && order != "DESC" {
		order = "ASC"
	}
	loadRatings, err := strconv.ParseBool(r.URL.Query().Get("ratings"))
	if err != nil {
		loadRatings = true
	}
	tweets := models.SearchTweets(mux.Vars(r)["needle"], uint(limit), order, loadRatings)
	json.NewEncoder(w).Encode(tweets)
}
