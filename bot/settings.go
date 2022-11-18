package bot

import (
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/samber/lo"
)

type toggleOptions struct {
	key                           string
	defaultValue                  bool
	enableMessage, disableMessage string
}

func (b *Bot) onToggle(opts toggleOptions) messageHandler {
	return func(message *telego.Message) error {
		return b.db.Update(func(tx *bolt.Tx) error {
			bk, err := tx.CreateBucketIfNotExists([]byte(opts.key))
			if err != nil {
				return err
			}

			id := []byte(strconv.FormatInt(message.From.ID, 10))
			v := bk.Get(id)

			def := []byte{0}
			if opts.defaultValue {
				def = []byte{1}
			}

			if v == nil {
				v = def
			}
			if v[0] == 0 {
				v = []byte{1}
			} else {
				v = []byte{0}
			}

			msg := opts.enableMessage
			if v[0] == 0 {
				msg = opts.disableMessage
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
}

func (b *Bot) getToggleValue(key string, id int64, defaultValue ...bool) (value bool) {
	b.db.View(func(tx *bolt.Tx) error {
		if bk := tx.Bucket([]byte(key)); bk != nil {
			id := []byte(strconv.FormatInt(id, 10))
			if v := bk.Get(id); len(v) > 0 {
				value = v[0] == 1
			} else if len(defaultValue) > 0 {
				value = defaultValue[0]
			}
		} else if len(defaultValue) > 0 {
			value = defaultValue[0]
		}
		return nil
	})
	return
}
