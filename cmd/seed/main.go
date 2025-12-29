package main

import (
	"log"
	"os"
	"time"

	"koalbot_api/internal/db"
	"koalbot_api/internal/seed"
)

func main() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN is required")
	}

	database, err := db.Open(dsn, 2, 2, 5*time.Minute)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	if err := seed.Users(database); err != nil {
		log.Fatal(err)
	}
}
