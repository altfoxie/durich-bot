package bot

import (
	"context"

	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
)

func (b *Bot) onStart(ctx context.Context, message *tg.Message, builder *message.RequestBuilder) error {
	_, err := builder.Text(ctx, "привет отправь текст я найду картинку и сделаю ржаку")
	return err
}
