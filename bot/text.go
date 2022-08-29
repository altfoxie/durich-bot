package bot

import (
	"bytes"
	"errors"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"

	"github.com/altfoxie/durich-bot/idraw"
	"github.com/altfoxie/durich-bot/vkapi"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/samber/lo"
)

func (b *Bot) onText(message *telego.Message) error {
	msg, err := b.SendMessage(
		tu.Message(tu.ID(message.Chat.ID), "🔎 ищем картинковое..."),
	)
	if err != nil {
		return err
	}

	edit := func(text string) error {
		return lo.T2(b.EditMessageText(&telego.EditMessageTextParams{
			ChatID:    tu.ID(msg.Chat.ID),
			MessageID: msg.MessageID,
			Text:      text,
		})).B
	}

	photo, err := vkapi.SearchRandomPhoto(message.Text)
	if err != nil {
		if err := edit("🤯 не нашел картинку ж есть"); err != nil {
			return err
		}
		return err
	}

	best := photo.BestSize()
	if best == nil {
		if err := edit("🤯 не нашел ссылку ж есть"); err != nil {
			return err
		}
		return errors.New("no best size")
	}

	if err = edit("💾 подожди ща скачаем..."); err != nil {
		return err
	}

	resp, err := http.Get(best.URL)
	if err != nil {
		if err := edit("🤯 не скачалось ж есть"); err != nil {
			return err
		}
		return err
	}

	img, err := jpeg.Decode(resp.Body)
	if err != nil {
		if err := edit("🤯 не отдекодилось ж есть"); err != nil {
			return err
		}
		return err
	}

	if err = edit("😸 ща прикол делаю..."); err != nil {
		return err
	}

	meme, err := idraw.MakeMeme(img, message.Text)
	if err != nil {
		if err := edit("🤯 мем не получился ж есть"); err != nil {
			return err
		}
		return err
	}

	buf := new(bytes.Buffer)
	if err := png.Encode(buf, meme); err != nil {
		return err
	}

	if err = b.DeleteMessage(&telego.DeleteMessageParams{
		ChatID:    tu.ID(msg.Chat.ID),
		MessageID: msg.MessageID,
	}); err != nil {
		log.Println("DeleteMessage error:", err)
	}

	return lo.T2(b.SendPhoto(
		tu.Photo(
			tu.ID(message.Chat.ID),
			tu.File(tu.NameReader(buf, "meme.png")),
		).WithReplyToMessageID(message.MessageID),
	)).B
}
