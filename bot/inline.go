package bot

import (
	"errors"
	"os"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func (b *Bot) onInlineQuery(query *telego.InlineQuery) error {
	answer := tu.InlineQuery(query.ID).WithIsPersonal().WithCacheTime(-1)

	if query.Query == "" {
		return b.AnswerInlineQuery(
			answer.WithSwitchPmText("напиши прекол ж есть").WithSwitchPmParameter("lol"),
		)
	}

	meme, err := makeMeme(query.Query)
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
		if err := b.AnswerInlineQuery(answer.WithSwitchPmText(errText).WithSwitchPmParameter("lol")); err != nil {
			return err
		}
		return err
	}

	msg, err := b.SendPhoto(
		tu.Photo(tu.Username(os.Getenv("CHANNEL_ID")), tu.File(tu.NameReader(meme, "meme.png"))),
	)
	if err != nil {
		return err
	}
	if len(msg.Photo) == 0 {
		return errors.New("no photo in message")
	}

	return b.AnswerInlineQuery(
		answer.WithResults(tu.ResultCachedPhoto("meme", msg.Photo[0].FileID)),
	)
}
