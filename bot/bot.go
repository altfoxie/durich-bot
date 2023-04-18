package bot

import (
	"context"
	"errors"
	"os"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/gotd/contrib/bg"
	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
)

const (
	defaultAppID   = 6
	defaultAppHash = "eb06d4abfb49dc3eeb1aeb98ae0f581e"
)

// Workaround to reduce boilerplate.
var ctx = context.Background()

type Bot struct {
	token  string
	client *telegram.Client
	db     *bolt.DB

	stop       bg.StopFunc
	dispatcher tg.UpdateDispatcher

	self     *tg.User
	username string
}

func New(token string, db *bolt.DB) (*Bot, error) {
	id, hash := defaultAppID, defaultAppHash
	if idEnv, hashEnv := os.Getenv("APP_ID"), os.Getenv("APP_HASH"); idEnv != "" && hashEnv != "" {
		id, _ = strconv.Atoi(idEnv)
		hash = hashEnv
	}
	if id <= 0 || hash == "" {
		return nil, errors.New("invalid app id or hash")
	}

	dispatcher := tg.NewUpdateDispatcher()

	s := os.Getenv("SESSION")
	if s == "" {
		s = "session.json"
	}

	client := telegram.NewClient(id, hash, telegram.Options{
		SessionStorage: &session.FileStorage{Path: s},
		UpdateHandler:  dispatcher,
	})

	return &Bot{
		token:      token,
		client:     client,
		db:         db,
		dispatcher: dispatcher,
	}, nil
}

func (b *Bot) Start() error {
	stop, err := bg.Connect(b.client)
	if err != nil {
		return err
	}
	b.stop = stop

	auth, err := b.client.Auth().Bot(ctx, b.token)
	if err != nil {
		return err
	}
	b.self = auth.GetUser().(*tg.User)
	b.username = b.self.Username

	b.dispatcher.OnNewMessage(b.onNewMessage)
	b.dispatcher.OnBotInlineQuery(b.onInlineQuery)

	return nil
}
