package main

import (
	"encoding/json"
	"os"

	"go.uber.org/zap"
)

type UserAccount struct {
	ID         string   `json:"_id"`
	IsActive   bool     `json:"isActive"`
	Balance    string   `json:"balance"`
	Picture    string   `json:"picture"`
	Age        int      `json:"age"`
	EyeColor   string   `json:"eyeColor"`
	Name       string   `json:"name"`
	Gender     string   `json:"gender"`
	Company    string   `json:"company"`
	Email      string   `json:"email"`
	Phone      string   `json:"phone"`
	Address    string   `json:"address"`
	About      string   `json:"about"`
	Registered string   `json:"registered"`
	Latitude   float64  `json:"latitude"`
	Longitude  float64  `json:"longitude"`
	Tags       []string `json:"tags"`
}

type Db struct {
	data map[string]*UserAccount
}

func NewDb() *Db {
	file, err := os.Open("data.json")
	if err != nil {
		logger.Fatal("can't read file", zap.Error(err))
	}
	defer file.Close()

	var decodedData []*UserAccount
	if err := json.NewDecoder(file).Decode(&decodedData); err != nil {
		logger.Fatal("can't load data from json", zap.Error(err))
	}
	db := Db{
		data: make(map[string]*UserAccount, len(decodedData)),
	}
	for _, r := range decodedData {
		db.data[r.ID] = r
	}
	logger.Info("db loaded", zap.Int("records", len(db.data)))
	return &db
}
