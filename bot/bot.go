package bot

import (
	"log"
	"os"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

type Bot struct {
	*telego.Bot
}

func New(token string) (*Bot, error) {
	bot, err := telego.NewBot(
		token,
		telego.WithDefaultLogger(os.Getenv("TELEGO_DEBUG") != "", true),
	)
	if err != nil {
		return nil, err
	}
	return &Bot{bot}, nil
}

func (b *Bot) Start() error {
	updates, err := b.UpdatesViaLongPulling(nil)
	if err != nil {
		return err
	}

	bh, err := th.NewBotHandler(b.Bot, updates)
	if err != nil {
		return err
	}

	bh.Handle(wrapMessageHandler(b.onStart), th.CommandEqual("start"))
	bh.Handle(wrapMessageHandler(b.onText), th.AnyMessageWithText())

	bh.Handle(wrapInlineQueryHandler(b.onInlineQuery), th.AnyInlineQuery())

	bh.Start()
	return nil
}

type messageHandler = func(message *telego.Message) error

func wrapMessageHandler(h messageHandler) th.Handler {
	return func(_ *telego.Bot, update telego.Update) {
		if err := h(update.Message); err != nil {
			log.Println("Message handler error:", err)
		}
	}
}

type inlineQueryHandler = func(query *telego.InlineQuery) error

func wrapInlineQueryHandler(h inlineQueryHandler) th.Handler {
	return func(_ *telego.Bot, update telego.Update) {
		if err := h(update.InlineQuery); err != nil {
			log.Println("Inline query handler error:", err)
		}
	}
}
