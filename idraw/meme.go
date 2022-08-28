package idraw

import (
	"image"
	"image/draw"

	"github.com/nfnt/resize"
)

func MakeMeme(img image.Image, text string) (image.Image, error) {
	textImg, err := DrawText(text, &TextOptions{
		Size:     72,
		MaxWidth: 600,
	})
	if err != nil {
		return nil, err
	}

	meme := image.NewRGBA(image.Rect(0, 0, 800, 600))
	draw.Draw(meme, meme.Bounds(), image.Black, image.ZP, draw.Src)

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
