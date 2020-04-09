package models

import (
	"errors"
	"time"
)

type APIKeyLog struct {
	ID         uint      `json:"-" gorm:"primary_key;not null"`
	APIKeyID   uint      `json:"-" gorm:"not null"`
	Path       string    `json:"path" gorm:"not null"`
	Method     string    `json:"method" gorm:"not null"`
	Authorized bool      `json:"authorized"`
	Date       time.Time `json:"date" gorm:"not null"`
}

func CreateAPIKeyLog(apiKeyID uint, path string, method string, authorized bool) error {
	log := &APIKeyLog{}
	log.APIKeyID = apiKeyID
	log.Path = path
	log.Method = method
	log.Authorized = authorized
	log.Date = time.Now()
	if GetDB().Create(&log).Error != nil {
		return errors.New("log could not be created due to a unexpected error")
	}
	return nil
}
