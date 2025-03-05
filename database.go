package main

import (
	"log"

	"gorm.io/driver/sqlite" // Sqlite driver based on CGO
	"gorm.io/gorm"

	"github.com/hex4coder/go-url-shortener/models"
)

var DB *gorm.DB

func InitDB() {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("error in database init : %v", err)
		return
	}

	err = db.AutoMigrate(&models.DataLink{}, &models.ShortLink{})
	if err != nil {
		log.Fatalf("failed to create models in database : %v", err)
		return
	}
	DB = db
}
