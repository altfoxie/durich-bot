package idraw

import (
	"bytes"
	_ "embed"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"image/draw"
	"image/png"
)

//go:embed banner.png
var banner []byte

func MakeMeme(img image.Image, text string) (image.Image, error) {
	// Text
	textImg, err := DrawText(text, &TextOptions{
		Size:     72,
		MaxWidth: 600,
	})
	if err != nil {
		return nil, err
	}

	// Watermark
	watermark, err := DrawText("@durich_bot", &TextOptions{
		Color: color.Alpha{A: 100},
		Size:  16,
	})
	if err != nil {
		return nil, err
	}

	meme := image.NewRGBA(image.Rect(0, 0, 800, 800))
	draw.Draw(meme, meme.Bounds(), image.Black, image.ZP, draw.Src)

	// 1XBET BANNER XDD
	banner, _ := png.Decode(bytes.NewReader(banner))
	draw.Draw(meme, image.Rect(0, 600, 800, 800), banner, image.ZP, draw.Over)

	// Base image
	img = resize.Resize(600, 400, img, resize.Bicubic)
	x0 := (meme.Bounds().Dx() - img.Bounds().Dx()) / 2
	drawRect := image.Rect(x0, 50, x0+img.Bounds().Dx(), 50+img.Bounds().Max.Y)
	draw.Draw(
		meme,
		drawRect,
		img,
		image.ZP,
		draw.Over,
	)

	// Base image borders
	DrawBorders(meme, image.Rect(
		drawRect.Min.X-3,
		drawRect.Min.Y-3,
		drawRect.Max.X+3,
		drawRect.Max.Y+3,
	), &BordersOptions{Size: 3})

	// Watermark
	x0 = drawRect.Max.X - watermark.Bounds().Dx()
	draw.Draw(meme, image.Rect(
		x0,
		drawRect.Min.Y-16-16,
		x0+watermark.Bounds().Dx(),
		drawRect.Min.Y-16-16+watermark.Bounds().Dy(),
	), watermark, image.ZP, draw.Over)

	// Text
	x0 = (meme.Bounds().Dx() - textImg.Bounds().Dx()) / 2
	drawRect = image.Rect(
		x0,
		drawRect.Max.Y+24,
		x0+textImg.Bounds().Dx(),
		drawRect.Max.Y+24+textImg.Bounds().Dy(),
	)
	draw.Draw(
		meme,
		drawRect,
		textImg,
		image.ZP,
		draw.Over,
	)

	return meme, nil
}
