package bot

import (
	"context"
	"errors"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
)

type toggleOptions struct {
	key                           string
	defaultValue                  bool
	enableMessage, disableMessage string
}

func (b *Bot) onToggle(ctx context.Context, message *tg.Message, builder *message.RequestBuilder, opts toggleOptions) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		bk, err := tx.CreateBucketIfNotExists([]byte(opts.key))
		if err != nil {
			return err
		}

		user, ok := message.GetPeerID().(*tg.PeerUser)
		if !ok {
			return errors.New("unexpected peer type")
		}

		id := []byte(strconv.FormatInt(user.UserID, 10))
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

		_, err = builder.Text(ctx, msg)
		return err
	})
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
