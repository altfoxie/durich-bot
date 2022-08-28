package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/altfoxie/durich-bot/bot"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	b, err := bot.New(os.Getenv("TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	if err = b.Start(); err != nil {
		log.Fatal(err)
	}
}
