package vkapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"
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

func SearchPhoto(query string, random bool) (*Photo, error) {
	resp, err := http.Get(
		"https://api.vk.com/method/photos.search?" + url.Values{
			"access_token": {os.Getenv("VK_TOKEN")},
			"q":            {query},
			"sort":         {"0"},
			// "count":        {"5"},
			"v": {"5.123"},
		}.Encode(),
	)
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

	if !random {
		for _, p := range r.Response.Items {
			if s := p.BestSize(); s != nil {
				return &p, nil
			}
		}
	}

	unique := make(map[string]Photo, len(r.Response.Items)) // URL -> Photo
	for _, p := range r.Response.Items {
		if s := p.BestSize(); s != nil {
			unique[s.URL] = p
		}
	}

	var list []Photo
	for _, p := range unique {
		list = append(list, p)
	}
	if len(list) == 0 {
		// Old method
		list = r.Response.Items
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	return &list[rnd.Intn(len(list))], nil
}
