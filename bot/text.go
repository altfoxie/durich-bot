package bot

import (
	"bytes"
	"errors"
	"image/jpeg"
	"image/png"
	"net/http"

	"github.com/altfoxie/durich-bot/idraw"
	"github.com/altfoxie/durich-bot/vkapi"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/samber/lo"
)

func (b *Bot) onText(message *telego.Message) error {
	photo, err := vkapi.SearchPhoto(message.Text)
	if err != nil {
		return err
	}

	best := photo.BestSize()
	if best == nil {
		return errors.New("no best size")
	}

	resp, err := http.Get(best.URL)
	if err != nil {
		return err
	}

	img, err := jpeg.Decode(resp.Body)
	if err != nil {
		return err
	}

	meme, err := idraw.MakeMeme(img, message.Text)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	if err := png.Encode(buf, meme); err != nil {
		return err
	}

	return lo.T2(
		b.SendPhoto(
			tu.Photo(
				tu.ID(message.Chat.ID),
				tu.File(tu.NameReader(buf, "meme.png")),
			),
		),
	).B
}
