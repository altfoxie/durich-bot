package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/altfoxie/durich-bot/bot"
	"github.com/boltdb/bolt"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	rand.Seed(time.Now().UnixNano())

	db, err := bolt.Open("bot.db", 0o600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	b, err := bot.New(
		os.Getenv("API_URL"),
		os.Getenv("TOKEN"),
		db,
	)
	if err != nil {
		log.Fatal(err)
	}

	if err = b.Start(); err != nil {
		log.Fatal(err)
	}

	select {}
}
