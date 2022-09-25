package bot

import (
	"bytes"
	"errors"
	"fmt"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/altfoxie/durich-bot/idraw"
	"github.com/altfoxie/durich-bot/vkapi"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/samber/lo"
)

var (
	errImageNotFound    = errors.New("image not found")
	errBestSizeNotFound = errors.New("best size not found")
	errImageGet         = errors.New("image get error")
	errDecode           = errors.New("image decode error")
	errMeme             = errors.New("make meme error")
	errEncode           = errors.New("image encode error")
)

type wrappedError struct {
	err          error
	recognizable error
}

func (e wrappedError) Error() string {
	if e.recognizable != nil {
		return e.recognizable.Error() + ": " + e.err.Error()
	}
	return e.err.Error()
}

func wrapError(err error, recognizable error) error {
	return wrappedError{err, recognizable}
}

func errorIs(err, target error) bool {
	if wrapped, ok := err.(wrappedError); ok {
		return errors.Is(wrapped.recognizable, target)
	}
	return errors.Is(err, target)
}

func (b *Bot) onMeme(message *telego.Message) error {
	msg, err := b.SendMessage(
		tu.Message(tu.ID(message.Chat.ID), "😸 ща прикол сделаю....").
			WithReplyToMessageID(message.MessageID),
	)
	if err != nil {
		return err
	}

	var meme io.Reader
	if len(message.Photo) > 0 {
		meme, err = b.makeMemeFromPhoto(
			message.Photo[len(message.Photo)-1],
			message.Caption,
		)
	} else {
		meme, err = makeMeme(message.Text)
	}
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

func makeMeme(query string) (io.Reader, error) {
	photo, err := vkapi.SearchRandomPhoto(strings.Split(query, "\n")[0])
	if err != nil {
		return nil, wrapError(err, errImageNotFound)
	}

	best := photo.BestSize()
	if best == nil {
		return nil, wrapError(err, errBestSizeNotFound)
	}

	return makeMemeFromURL(best.URL, query)
}

func (b *Bot) makeMemeFromPhoto(
	photo telego.PhotoSize,
	text string,
) (io.Reader, error) {
	file, err := b.GetFile(&telego.GetFileParams{
		FileID: photo.FileID,
	})
	if err != nil {
		return nil, wrapError(err, errImageGet)
	}

	return makeMemeFromURL(
		fmt.Sprintf(
			"https://api.telegram.org/file/bot%s/%s",
			b.Token(),
			file.FilePath,
		),
		text,
	)
}

func makeMemeFromURL(url, text string) (io.Reader, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, wrapError(err, errImageGet)
	}

	img, err := jpeg.Decode(resp.Body)
	if err != nil {
		return nil, wrapError(err, errDecode)
	}

	meme, err := idraw.MakeMeme(img, text)
	if err != nil {
		return nil, wrapError(err, errMeme)
	}

	buf := new(bytes.Buffer)
	if err := png.Encode(buf, meme); err != nil {
		return nil, wrapError(err, errEncode)
	}

	return buf, nil
}
