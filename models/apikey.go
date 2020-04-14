package models

import (
	"encoding/base64"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type AccessLevel uint

type APIKey struct {
	ID          uint        `json:"-" gorm:"primary_key;not null"`
	Name        string      `json:"name" gorm:"not null"`
	Key         string      `json:"-" gorm:"not null"`
	AccessLevel AccessLevel `json:"access_level" gorm:"not null"`
	Creation    time.Time   `json:"creation" gorm:"not null"`
	LastUse     time.Time   `json:"last_use" gorm:"not null"`
	Logs        []APIKeyLog `json:"logs,omitempty" gorm:"foreignkey:ApiKeyID;association_autoupdate:false;association_autocreate:false"`
}

const (
	AccessNothing   AccessLevel = 0
	AccessReadOnly  AccessLevel = 1
	AccessReadWrite AccessLevel = 2
)

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func (k *APIKey) CanRead() bool {
	return k.AccessLevel >= AccessReadOnly
}

func (k *APIKey) CanWrite() bool {
	return k.AccessLevel >= AccessReadWrite
}

func (k *APIKey) valid(token string) bool {
	return bcrypt.CompareHashAndPassword([]byte(k.Key), []byte(token)) == nil
}

func getById(id uint) APIKey {
	apiKey := &APIKey{}
	GetDB().First(apiKey, id)
	return *apiKey
}

func (k *APIKey) LoadLogs(limit uint) error {
	if GetDB().Model(&k).Limit(limit).Order("id DESC").Related(&k.Logs).Error != nil {
		return errors.New("could not load logs")
	}
	return nil
}

func (k *APIKey) RegisterUse(path string, method string, authorized bool) error {
	err := CreateAPIKeyLog(k.ID, path, method, authorized)
	if err != nil {
		return err
	}
	if GetDB().Model(&k).Update("last_use", time.Now()) != nil {
		return errors.New("api key last use date could not be updated")
	}
	return nil
}

func GetAPIKey(token string) (*APIKey, error) {
	split := strings.Split(token, ".")
	if len(split) != 2 {
		return nil, errors.New("invalid token")
	}
	id, err := base64.StdEncoding.DecodeString(split[0])
	if err != nil {
		return nil, errors.New("invalid token")
	}
	idInt, err := strconv.Atoi(string(id))
	if err != nil {
		return nil, errors.New("invalid token")
	}
	apiKey := getById(uint(idInt))
	if !apiKey.valid(split[1]) {
		apiKey.Key = ""
		return nil, errors.New("invalid token")
	}
	return &apiKey, nil
}

func CreateAPIKey(name string) (string, error) {
	var err error
	var hashedKey []byte
	apiKey := &APIKey{}
	key := generateKey()
	cost := 10
	apiKey.AccessLevel = AccessReadOnly
	hashedKey, err = bcrypt.GenerateFromPassword([]byte(key), cost)
	if err != nil {
		return "", errors.New("api key could not be created due to a unexpected error")
	}
	apiKey.Name = name
	apiKey.Key = string(hashedKey)
	apiKey.Creation = time.Now()
	if GetDB().Create(&apiKey).Error != nil {
		return "", errors.New("api key could not be created due to a unexpected error")
	}
	apiKey.RegisterUse(os.Getenv("path")+"apikeys/new", http.MethodGet, true)
	id := base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(int(apiKey.ID))))
	return id + "." + key, nil
}

func generateKey() string {
	length := 18
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+?!()&%$@-_"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
