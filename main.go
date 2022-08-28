package main

import (
	"log"
	"os"

	"github.com/altfoxie/durich-bot/bot"
)

func main() {
	b, err := bot.New(os.Getenv("TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	if err = b.Start(); err != nil {
		log.Fatal(err)
	}
}
