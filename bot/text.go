package bot

import (
	"log"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/samber/lo"
)

func (b *Bot) onText(message *telego.Message) error {
	msg, err := b.SendMessage(
		tu.Message(tu.ID(message.Chat.ID), "😸 ща прикол сделаю...."),
	)
	if err != nil {
		return err
	}

	meme, err := makeMeme(message.Text)
	if err != nil {
		errText := "🤯 неизвестная ошибка ж есть"
		switch {
		case errorIs(err, errImageNotFound):
			errText = "🤯 не нашел картинку ж есть"
		case errorIs(err, errBestSizeNotFound):
			errText = "🤯 не нашел ссылку ж есть"
		case errorIs(err, errImageGet):
			errText = "🤯 не скачалось ж есть"
		case errorIs(err, errDecode):
			errText = "🤯 не отдекодилось ж есть"
		case errorIs(err, errMeme):
			errText = "🤯 мем не получился ж есть"
		case errorIs(err, errEncode):
			errText = "🤯 не отдекодилось ж есть"
		}
		if _, err := b.EditMessageText(&telego.EditMessageTextParams{
			ChatID:    tu.ID(msg.Chat.ID),
			MessageID: msg.MessageID,
			Text:      errText,
		}); err != nil {
			return err
		}
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
			tu.File(tu.NameReader(meme, "meme.png")),
		).WithReplyToMessageID(message.MessageID),
	)).B
}
