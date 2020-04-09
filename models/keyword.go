package models

import (
	"errors"
	"github.com/jinzhu/gorm"
)

type Keyword struct {
	Text     string   `json:"text" gorm:"primary_key;not null"`
	Counter  uint     `json:"counter" gorm:"not null;"`
	Searches uint     `json:"searches" gorm:"not null"`
	Ratings  []Rating `json:"ratings,omitempty" gorm:"foreignkey:KeywordText"`
}

func (k *Keyword) incrementSearches() {
	GetDB().Model(&k).Update("searches", gorm.Expr("searches + ?", 1))
	k.Searches++
}

func (k *Keyword) incrementCounter() {
	GetDB().Model(&k).Update("counter", gorm.Expr("counter + ?", 1))
	k.Counter++
}

func (k *Keyword) Validate() error {
	if k.Text == "" {
		return errors.New("text field cannot be null")
	}
	return nil
}

func (k *Keyword) CreateOrUpdate() error {
	if err := k.Validate(); err != nil {
		return err
	}
	if GetDB().FirstOrCreate(&k).Error != nil {
		return errors.New("unexpected error while trying to create keyword")
	}
	k.incrementCounter()
	return nil
}

func GetKeyword(text string) *Keyword {
	keyword := &Keyword{}
	if GetDB().Where("text ILIKE ?", text).First(&keyword).Error != nil {
		return nil
	}
	ratings := make([]Rating, 0)
	GetDB().Select([]string{"rate", "tweet_id"}).Model(&keyword).Related(&ratings, "Ratings")
	for i, r := range ratings {
		tweet := &Tweet{}
		GetDB().Model(&r).Related(tweet)
		ratings[i].Tweet = tweet
		ratings[i].TweetID = 0
	}
	keyword.Ratings = ratings
	return keyword
}

func GetKeywords(limit uint, orderBy string, loadRatings bool) []*Keyword {
	if orderBy != "counter" && orderBy != "searches" {
		return nil
	}
	keywords := make([]*Keyword, 0)
	if GetDB().Order(orderBy+" desc").Limit(limit).Find(&keywords).Error != nil {
		return nil
	}
	if !loadRatings {
		return keywords
	}
	for _, k := range keywords {
		ratings := make([]Rating, 0)
		GetDB().Select([]string{"rate", "tweet_id"}).Model(&k).Related(&ratings, "Ratings")
		for i, r := range ratings {
			tweet := &Tweet{}
			GetDB().Model(&r).Related(tweet)
			ratings[i].Tweet = tweet
			ratings[i].TweetID = 0
		}
		k.Ratings = ratings
	}
	return keywords
}

func SearchKeywords(needle string, limit uint, loadRatings bool) []*Keyword {
	keywords := make([]*Keyword, 0)
	if GetDB().Where("text ILIKE ?", "%"+needle+"%").Limit(limit).Find(&keywords).Error != nil {
		return nil
	}
	if !loadRatings {
		return keywords
	}
	for _, k := range keywords {
		k.incrementSearches()
		ratings := make([]Rating, 0)
		GetDB().Select([]string{"rate", "tweet_id"}).Model(&k).Related(&ratings, "Ratings")
		for i, r := range ratings {
			tweet := &Tweet{}
			GetDB().Model(&r).Related(tweet)
			ratings[i].Tweet = tweet
			ratings[i].TweetID = 0
		}
		k.Ratings = ratings
	}
	return keywords
}
