package models

import (
	"errors"
	"github.com/lib/pq"
	"strings"
	"time"
)

type Tweet struct {
	ID              uint           `json:"id" gorm:"primary_key;not null"`
	Title           string         `json:"title" gorm:"not null"`
	Image           string         `json:"image" gorm:"type:varchar(2083);not null"`
	ArchillectTweet uint64         `json:"archillect_tweet" gorm:"not null"`
	Date            time.Time      `json:"date" gorm:"not null"`
	Pages           pq.StringArray `json:"pages" gorm:"type:varchar(2083)[];not null"`
	Matches         pq.StringArray `json:"matches" gorm:"type:varchar(2083)[];not null"`
	PartialMatches  pq.StringArray `json:"partial_matches" gorm:"type:varchar(2083)[];not null"`
	NumKeywords     uint           `json:"num_keywords"`
	IsFirst         bool           `json:"first" gorm:"-"`
	IsLast          bool           `json:"last" gorm:"-"`
	Ratings         []Rating       `json:"ratings,omitempty"`
}

func (t *Tweet) validate() error {
	if t.Pages == nil {
		return errors.New("pages field cannot be null")
	}
	if t.Matches == nil {
		return errors.New("matches field cannot be null")
	}
	if t.PartialMatches == nil {
		return errors.New("partial_matches field cannot be null")
	}
	if t.Ratings == nil {
		return errors.New("ratings field cannot be null")
	}
	if t.Image == "" {
		return errors.New("image field cannot be empty")
	}
	if t.ArchillectTweet == 0 {
		return errors.New("archillect_tweet field cannot be empty")
	}
	return nil
}

func (t *Tweet) Create() error {
	if err := t.validate(); err != nil {
		return err
	}
	if GetDB().Create(t).Error != nil {
		return errors.New("unexpected error while trying to create tweet")
	}
	return nil
}

func GetTweet(id uint) *Tweet {
	tweet := &Tweet{}
	if GetDB().First(&tweet, id).Error != nil {
		return nil
	}
	ratings := make([]Rating, 0)
	GetDB().Select([]string{"rate", "keyword_text"}).Model(&tweet).Related(&ratings)
	for i, r := range ratings {
		keyword := &Keyword{}
		GetDB().Model(&r).Related(keyword, "KeywordText")
		ratings[i].Keyword = keyword
		ratings[i].KeywordText = ""
	}
	last := &Tweet{}
	GetDB().Last(&last)
	tweet.Ratings = ratings
	tweet.IsFirst = tweet.ID == 1
	tweet.IsLast = tweet.ID == last.ID
	return tweet
}

func GetTweets(limit uint, order string, loadRatings bool) []*Tweet {
	order = strings.ToUpper(order)
	if order != "ASC" && order != "DESC" {
		return nil
	}
	tweets := make([]*Tweet, 0)
	if GetDB().Limit(limit).Order("id "+order).Find(&tweets).Error != nil {
		return nil
	}
	last := &Tweet{}
	GetDB().Last(&last)
	for _, t := range tweets {
		t.IsFirst = t.ID == 1
		t.IsLast = t.ID == last.ID
		if loadRatings {
			ratings := make([]Rating, 0)
			GetDB().Select([]string{"rate", "keyword_text"}).Model(&t).Related(&ratings)
			for i, r := range ratings {
				keyword := &Keyword{}
				GetDB().Model(&r).Related(keyword, "KeywordText")
				ratings[i].Keyword = keyword
				ratings[i].KeywordText = ""
			}
			t.Ratings = ratings
		}
	}
	return tweets
}

func SearchTweets(needle string, limit uint, order string, loadRatings bool) []*Tweet {
	order = strings.ToUpper(order)
	if order != "ASC" && order != "DESC" {
		return nil
	}
	tweets := make([]*Tweet, 0)
	_, err := GetDB().Table("tweets").Select("*").
		Order("id "+order).
		Limit(limit).
		Where("EXISTS (SELECT 1 FROM ratings WHERE ratings.keyword_text ILIKE ? AND ratings.tweet_id = tweets.id)", "%"+needle+"%").
		Find(&tweets).Rows()
	if err != nil {
		return nil
	}
	last := &Tweet{}
	GetDB().Last(&last)
	for _, tweet := range tweets {
		tweet.IsFirst = tweet.ID == 1
		tweet.IsLast = tweet.ID == last.ID
		if loadRatings {
			ratings := make([]Rating, 0)
			GetDB().Select([]string{"rate", "keyword_text"}).Model(&tweet).Related(&ratings)
			for i, r := range ratings {
				keyword := &Keyword{}
				GetDB().Model(&r).Related(keyword, "KeywordText")
				ratings[i].Keyword = keyword
				ratings[i].KeywordText = ""
			}
			tweet.Ratings = ratings
		}
	}
	return tweets
}
