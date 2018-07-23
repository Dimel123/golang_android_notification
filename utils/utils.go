package utils

import (
	"net/http"
	"encoding/json"
	"testNotification/models"
)

func ParseToken(r *http.Request) models.Token {
	decoder := json.NewDecoder(r.Body)
	token := models.Token{}
	err := decoder.Decode(&token)
	if err != nil {
		panic(err)
	}
	return token
}

func ParseNotification(r *http.Request) models.Notification {
	decoder := json.NewDecoder(r.Body)
	notification := models.Notification{}
	err := decoder.Decode(&notification)
	if err != nil {
		panic(err)
	}
	return notification
}