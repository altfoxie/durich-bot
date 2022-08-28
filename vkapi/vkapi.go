package vkapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
)

type response[T any] struct {
	Response *T        `json:"response"`
	Error    *apiError `json:"error"`
}

type listResponse[T any] struct {
	Count int `json:"count"`
	Items []T `json:"items"`
}

type apiError struct {
	Code    int    `json:"error_code"`
	Message string `json:"error_msg"`
}

func (e *apiError) Error() string {
	return fmt.Sprintf("vkapi: %d - %s", e.Code, e.Message)
}

type Photo struct {
	ID      int         `json:"id"`
	OwnerID int         `json:"owner_id"`
	AlbumID int         `json:"album_id"`
	Date    int         `json:"date"`
	Sizes   []PhotoSize `json:"sizes"`
}

func (p *Photo) BestSize() *PhotoSize {
	var best *PhotoSize
	for _, s := range p.Sizes {
		if best == nil || s.Width > best.Width || s.Height > best.Height {
			best = &s
		}
	}
	return best
}

type PhotoSize struct {
	Type   string `json:"type"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	URL    string `json:"url"`
}

func SearchRandomPhoto(query string) (*Photo, error) {
	resp, err := http.Get("https://api.vk.com/method/photos.search?" + url.Values{
		"access_token": {os.Getenv("VK_TOKEN")},
		"q":            {query},
		// "count":        {"5"},
		"v": {"5.123"},
	}.Encode())
	if err != nil {
		return nil, err
	}

	var r response[listResponse[Photo]]
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}

	if r.Error != nil {
		return nil, r.Error
	}

	if len(r.Response.Items) == 0 {
		return nil, errors.New("vkapi: no photos found")
	}

	return &r.Response.Items[rand.Intn(len(r.Response.Items))], nil
}
