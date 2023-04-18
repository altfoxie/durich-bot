package bot

import (
	"context"
	"crypto/rand"
	"errors"
	"log"

	"github.com/gotd/td/telegram/message/inline"
	"github.com/gotd/td/telegram/message/markup"
	"github.com/gotd/td/telegram/uploader"
	"github.com/gotd/td/tg"
)

func (b *Bot) onInlineQuery(ctx context.Context, entities tg.Entities, update *tg.UpdateBotInlineQuery) error {
	builder := inline.New(b.client.API(), rand.Reader, update.QueryID).
		CacheTimeSeconds(-1).
		Private(true)

	if update.Query == "" {
		_, err := builder.SwitchPM("напиши прекол ж есть", "lol").Set(ctx)
		return err
	}

	reader, link, err := b.memeSearch(update.Query)
	_ = link
	if err == nil {
		reader, err = b.onMemeReader(update.Query, update.UserID, reader)
	}
	if err != nil {
		errText := "неизвестная ошибка ж есть"
		switch {
		case errorIs(err, errImageNotFound):
			errText = "не нашел картинку ж есть"
		case errorIs(err, errBestSizeNotFound):
			errText = "не нашел ссылку ж есть"
		case errorIs(err, errImageGet):
			errText = "не скачалось ж есть"
		case errorIs(err, errDecode):
			errText = "не отдекодилось ж есть"
		case errorIs(err, errMeme):
			errText = "мем не получился ж есть"
		case errorIs(err, errEncode):
			errText = "не отдекодилось ж есть"
		}
		if _, err := builder.SwitchPM(errText, "lol").Set(ctx); err != nil {
			return err
		}
		return err
	}

	file, err := uploader.NewUploader(b.client.API()).FromReader(ctx, "meme.png", reader)
	if err != nil {
		return err
	}

	media, err := b.client.API().MessagesUploadMedia(ctx, &tg.MessagesUploadMediaRequest{
		Peer: b.self.AsInputPeer(),
		Media: &tg.InputMediaUploadedPhoto{
			File: file,
		},
	})
	if err != nil {
		return err
	}

	photo, ok := media.(*tg.MessageMediaPhoto).Photo.(*tg.Photo)
	if !ok {
		return errors.New("unexpected inline media")
	}

	var replyMarkup tg.ReplyMarkupClass
	if link != "" && b.getToggleValue("link", update.UserID, true) {
		replyMarkup = markup.InlineRow(markup.URL("🔗 Ссылка", link))
	}

	result := inline.Photo(&tg.InputPhoto{
		ID:            photo.ID,
		AccessHash:    photo.AccessHash,
		FileReference: photo.FileReference,
	}, inline.ResultMessage(&tg.InputBotInlineMessageMediaAuto{
		ReplyMarkup: replyMarkup,
	}))

	_, err = builder.Gallery(true).Set(ctx, result)
	if err != nil {
		log.Println("Inline query error:", err)
	}
	return err
}
