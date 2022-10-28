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
	"strconv"
	"strings"

	"github.com/altfoxie/durich-bot/idraw"
	"github.com/altfoxie/durich-bot/vkapi"
	"github.com/boltdb/bolt"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/nfnt/resize"
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

	zhmyh := false
	b.db.View(func(tx *bolt.Tx) error {
		if bk := tx.Bucket([]byte("zhmyh")); bk != nil {
			id := []byte(strconv.FormatInt(message.From.ID, 10))
			if v := bk.Get(id); len(v) > 0 {
				zhmyh = v[0] == 1
			}
		}
		return nil
	})

	var meme io.Reader
	if len(message.Photo) > 0 {
		meme, err = b.makeMemeFromPhoto(
			message.Photo[len(message.Photo)-1],
			message.Caption,
			zhmyh,
		)
	} else {
		meme, err = makeMeme(message.Text, zhmyh)
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

func makeMeme(query string, zhmyh bool) (io.Reader, error) {
	photo, err := vkapi.SearchRandomPhoto(strings.Split(query, "\n")[0])
	if err != nil {
		return nil, wrapError(err, errImageNotFound)
	}

	best := photo.BestSize()
	if best == nil {
		return nil, wrapError(err, errBestSizeNotFound)
	}

	return makeMemeFromURL(best.URL, query, zhmyh)
}

func (b *Bot) makeMemeFromPhoto(
	photo telego.PhotoSize,
	text string,
	zhmyh bool,
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
		zhmyh,
	)
}

func makeMemeFromURL(url, text string, zhmyh bool) (io.Reader, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, wrapError(err, errImageGet)
	}

	img, err := jpeg.Decode(resp.Body)
	if err != nil {
		return nil, wrapError(err, errDecode)
	}

	if zhmyh {
		img = resize.Resize(600, 400, img, resize.Bilinear)
	}

	layers := strings.Split(text, "\n\n")
	fmt.Println(layers)

	for i, layer := range layers {
		lines := strings.SplitN(layer, "\n", 2)
		secondLine := ""
		if len(lines) > 1 {
			secondLine = lines[1]
		}

		if img, err = idraw.MakeMeme(img, lines[0], secondLine, i == len(layers)-1); err != nil {
			return nil, wrapError(err, errMeme)
		}
	}

	buf := new(bytes.Buffer)
	if err := png.Encode(buf, img); err != nil {
		return nil, wrapError(err, errEncode)
	}

	return buf, nil
}
