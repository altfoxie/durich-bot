package idraw

import (
	"image"
	"image/color"
	"image/draw"
	"strings"

	"github.com/nfnt/resize"
)

func MakeMeme(img image.Image, text string) (image.Image, error) {
	lines := strings.Split(text, "\n")
	firstLine, secondLine := lines[0], ""
	if len(lines) > 1 {
		secondLine = strings.Join(lines[1:], " ")
	}

	// Text (first line)
	firstLineImg, err := DrawText(firstLine, &TextOptions{
		Size:     72,
		MaxWidth: 600,
	})
	if err != nil {
		return nil, err
	}

	// Text (second line)
	var secondLineImg image.Image
	if secondLine != "" {
		secondLineImg, err = DrawText(secondLine, &TextOptions{
			Size:     42,
			MaxWidth: 600,
		})
		if err != nil {
			return nil, err
		}
	}

	// Watermark
	watermark, err := DrawText("@durich_bot", &TextOptions{
		Color: color.Alpha{A: 100},
		Size:  16,
	})
	if err != nil {
		return nil, err
	}

	scale := float64(600) / float64(img.Bounds().Dx())
	imgWidth, imgHeight := 600, int(float64(img.Bounds().Dy())*scale)
	height := imgHeight + 200
	if secondLineImg != nil {
		height += 50
	}

	meme := image.NewRGBA(image.Rect(0, 0, 800, height))
	draw.Draw(meme, meme.Bounds(), image.Black, image.ZP, draw.Src)

	// Base image
	img = resize.Resize(uint(imgWidth), uint(imgHeight), img, resize.Bicubic)
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

	// Text (first line)
	x0 = (meme.Bounds().Dx() - firstLineImg.Bounds().Dx()) / 2
	drawRect = image.Rect(
		x0,
		drawRect.Max.Y+24,
		x0+firstLineImg.Bounds().Dx(),
		drawRect.Max.Y+24+firstLineImg.Bounds().Dy(),
	)
	draw.Draw(
		meme,
		drawRect,
		firstLineImg,
		image.ZP,
		draw.Over,
	)

	// Text (second line)
	if secondLineImg != nil {
		x0 = (meme.Bounds().Dx() - secondLineImg.Bounds().Dx()) / 2
		drawRect = image.Rect(
			x0,
			drawRect.Max.Y,
			x0+secondLineImg.Bounds().Dx(),
			drawRect.Max.Y+secondLineImg.Bounds().Dy(),
		)
		draw.Draw(
			meme,
			drawRect,
			secondLineImg,
			image.ZP,
			draw.Over,
		)
	}

	return meme.SubImage(image.Rect(0, 0, 800, drawRect.Max.Y+50)), nil
}
