package db

import (
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error

	if _, err := os.Stat("data.db"); os.IsNotExist(err) {
		log.Println("📁 Creating new database file: data.db")
	}

	DB, err = gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Failed to open database: %v", err)
	}
	log.Println("📦 Connected to SQLite database (data.db)")
}
