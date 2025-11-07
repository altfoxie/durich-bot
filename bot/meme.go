package bot

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/altfoxie/durich-bot/vkapi"
	"github.com/gotd/td/telegram/message/markup"

	"github.com/altfoxie/durich-bot/idraw"
	"github.com/cognusion/go-utils/writeatbuffer"
	"github.com/gotd/td/telegram/downloader"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
	"github.com/nfnt/resize"
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

func (b *Bot) onMeme(ctx context.Context, msg *tg.Message, builder *message.RequestBuilder) error {
	var photo *tg.Photo
	if msg.Media != nil {
		mediaPhoto, ok := msg.Media.(*tg.MessageMediaPhoto)
		if !ok {
			return errors.New("unexpected media type")
		}

		if photo, ok = mediaPhoto.Photo.(*tg.Photo); !ok {
			return errors.New("unexpected media type")
		}
	} else if msg.Message == "" {
		return errors.New("empty message")
	}

	sentMsgUpdate, err := builder.ReplyMsg(msg).Text(ctx, "😸 ща прикол сделаю....")
	if err != nil {
		return err
	}

	sentMsg, ok := sentMsgUpdate.(*tg.UpdateShortSentMessage)
	if !ok {
		return errors.New("unexpected message type")
	}

	peer, ok := msg.GetPeerID().(*tg.PeerUser)
	if !ok {
		return errors.New("unexpected peer type")
	}

	var (
		reader     io.Reader
		buttonLink string
	)
	if photo != nil {
		buf := writeatbuffer.NewBuffer(make([]byte, 0, 1024*1024))
		if _, err = downloader.NewDownloader().Download(b.client.API(), &tg.InputPhotoFileLocation{
			ID:            photo.ID,
			AccessHash:    photo.AccessHash,
			FileReference: photo.FileReference,
			ThumbSize:     "x",
		}).Parallel(ctx, buf); err != nil {
			return wrapError(err, errImageGet)
		}

		reader = bytes.NewReader(buf.Bytes())
	} else {
		reader, buttonLink, err = b.memeSearch(msg.Message, peer.UserID)
	}

	if err == nil {
		reader, err = b.onMemeReader(msg.Message, peer.UserID, reader)
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
		if _, err := builder.Edit(sentMsg.ID).Text(ctx, errText); err != nil {
			return err
		}
		return err
	}

	if _, err = b.client.API().MessagesDeleteMessages(ctx, &tg.MessagesDeleteMessagesRequest{
		ID:     []int{sentMsg.ID},
		Revoke: true,
	}); err != nil {
		log.Println("delete message error:", err)
	}

	if buttonLink != "" && b.getToggleValue("link", peer.UserID, defaultLink) {
		*builder = message.RequestBuilder{
			Builder: *builder.Markup(markup.InlineRow(markup.URL("🔗 Ссылка", buttonLink))),
		}
	}

	_, err = builder.Upload(message.FromReader("meme.png", reader)).Photo(ctx)
	return err
}

func (b *Bot) memeSearch(text string, userID int64) (io.Reader, string, error) {
	photo, err := vkapi.SearchPhoto(strings.Split(text, "\n")[0], !b.getToggleValue("last", userID, defaultLast))
	if err != nil {
		return nil, "", wrapError(err, errImageNotFound)
	}

	best := photo.BestSize()
	if best == nil {
		return nil, "", wrapError(err, errBestSizeNotFound)
	}

	resp, err := http.Get(best.URL)
	if err != nil {
		return nil, "", wrapError(err, errImageGet)
	}

	return resp.Body, fmt.Sprintf("https://vk.com/photo%d_%d", photo.OwnerID, photo.ID), nil
}

func (b *Bot) onMemeReader(text string, userID int64, reader io.Reader) (io.Reader, error) {
	img, err := jpeg.Decode(reader)
	if err != nil {
		return nil, wrapError(err, errDecode)
	}

	zhmyh := b.getToggleValue("zhmyh", userID, defaultZhmyh)
	if zhmyh {
		img = resize.Resize(600, 400, img, resize.Bilinear)
	}

	layers := strings.Split(text, "\n\n")
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
