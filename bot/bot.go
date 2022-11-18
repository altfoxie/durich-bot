package bot

import (
	"log"
	"os"

	"github.com/boltdb/bolt"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

type Bot struct {
	*telego.Bot
	db *bolt.DB
}

func New(token string, db *bolt.DB) (*Bot, error) {
	bot, err := telego.NewBot(
		token,
		telego.WithDefaultLogger(os.Getenv("TELEGO_DEBUG") != "", true),
	)
	if err != nil {
		return nil, err
	}
	return &Bot{bot, db}, nil
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
	bh.Handle(wrapMessageHandler(b.onToggle(toggleOptions{
		key:            "zhmyh",
		enableMessage:  "теперь ты жмыхаешь картинки",
		disableMessage: "больше ты не жмыхаешь картинки",
	})), th.CommandEqual("zhmyh"))
	bh.Handle(wrapMessageHandler(b.onToggle(toggleOptions{
		key:            "link",
		defaultValue:   true,
		enableMessage:  "теперь я буду отправлять ссылки на картинки",
		disableMessage: "больше я не буду отправлять ссылки на картинки",
	})), th.CommandEqual("link"))
	bh.Handle(wrapMessageHandler(b.onMeme), func(update telego.Update) bool {
		return update.Message != nil && len(update.Message.Photo) > 0
	})
	bh.Handle(wrapMessageHandler(b.onMeme), th.AnyMessageWithText())

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
