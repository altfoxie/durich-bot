package bot

import (
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"

	"github.com/samber/lo"
)

func (b *Bot) onStart(message *telego.Message) error {
	return lo.T2(
		b.SendMessage(
			tu.Message(
				tu.ID(message.Chat.ID),
				"привет отправь текст я найду картинку и сделаю ржаку",
			),
		),
	).B
}
