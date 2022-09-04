package bot

import (
	"bytes"
	"errors"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"strings"

	"github.com/altfoxie/durich-bot/idraw"
	"github.com/altfoxie/durich-bot/vkapi"
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

func makeMeme(query string) (io.Reader, error) {
	photo, err := vkapi.SearchRandomPhoto(strings.Split(query, "\n")[0])
	if err != nil {
		return nil, wrapError(err, errImageNotFound)
	}

	best := photo.BestSize()
	if best == nil {
		return nil, wrapError(err, errBestSizeNotFound)
	}

	resp, err := http.Get(best.URL)
	if err != nil {
		return nil, wrapError(err, errImageGet)
	}

	img, err := jpeg.Decode(resp.Body)
	if err != nil {
		return nil, wrapError(err, errDecode)
	}

	meme, err := idraw.MakeMeme(img, query)
	if err != nil {
		return nil, wrapError(err, errMeme)
	}

	buf := new(bytes.Buffer)
	if err := png.Encode(buf, meme); err != nil {
		return nil, wrapError(err, errEncode)
	}

	return buf, nil
}
