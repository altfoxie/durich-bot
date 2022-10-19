package bot

import (
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/samber/lo"
)

func (b *Bot) onZhmyh(message *telego.Message) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		bk, err := tx.CreateBucketIfNotExists([]byte("zhmyh"))
		if err != nil {
			return err
		}

		id := []byte(strconv.FormatInt(message.From.ID, 10))
		v := bk.Get(id)
		if v == nil {
			v = []byte{0}
		}

		msg := "теперь ты жмыхаешь картинки"
		if v[0] == 0 {
			v[0] = 1
		} else {
			v[0] = 0
			msg = "больше ты не жмыхаешь картинки"
		}

		if err = bk.Put(id, v); err != nil {
			return err
		}

		return lo.T2(
			b.SendMessage(
				tu.Message(
					tu.ID(message.Chat.ID),
					msg,
				),
			),
		).B
	})
}
