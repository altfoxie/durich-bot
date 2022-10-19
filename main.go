package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/altfoxie/durich-bot/bot"
	"github.com/boltdb/bolt"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	db, err := bolt.Open("bot.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	b, err := bot.New(os.Getenv("TOKEN"), db)
	if err != nil {
		log.Fatal(err)
	}

	if err = b.Start(); err != nil {
		log.Fatal(err)
	}
}
