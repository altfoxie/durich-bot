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

	p := os.Getenv("DB")
	if p == "" {
		p = "bot.db"
	}

	db, err := bolt.Open(p, 0o600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	b, err := bot.New(
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
