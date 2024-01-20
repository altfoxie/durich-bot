package bot

import (
	"context"
	"log"
	"strings"

	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
)

func (b *Bot) onNewMessage(ctx context.Context, entities tg.Entities, update *tg.UpdateNewMessage) error {
	m, ok := update.Message.(*tg.Message)
	if !ok || m.Out {
		return nil
	}

	longCommand, _, _ := strings.Cut(m.Message, " ")
	sender := message.NewSender(b.client.API())
	builder := sender.Answer(entities, update)

	var err error
	if strings.HasPrefix(longCommand, "/") {
		command, botName, _ := strings.Cut(longCommand, "@")
		if botName != "" && !strings.EqualFold(b.username, botName) {
			return nil
		}

		switch {
		case strings.EqualFold(command, "/start"):
			err = b.onStart(ctx, m, builder)
		case strings.EqualFold(command, "/zhmyh"):
			err = b.onToggle(ctx, m, builder, toggleOptions{
				key:            "zhmyh",
				defaultValue:   true,
				enableMessage:  "теперь ты жмыхаешь картинки",
				disableMessage: "больше ты не жмыхаешь картинки",
			})
		case strings.EqualFold(command, "/link"):
			err = b.onToggle(ctx, m, builder, toggleOptions{
				key:            "link",
				defaultValue:   true,
				enableMessage:  "теперь я буду отправлять ссылки на картинки",
				disableMessage: "больше я не буду отправлять ссылки на картинки",
			})
		case strings.EqualFold(command, "/last"):
			err = b.onToggle(ctx, m, builder, toggleOptions{
				key:            "last",
				defaultValue:   true,
				enableMessage:  "теперь я буду использовать последнюю картинку из вкалтакте",
				disableMessage: "больше я не буду использовать последнюю картинку из вкалтакте. теперь тока случайную...",
			})
		}
	} else {
		err = b.onMeme(ctx, m, builder)
	}

	if err != nil {
		log.Println("Message handler error:", err)
	}
	return err
}
