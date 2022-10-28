package idraw

import (
	"image"
	"image/color"
	"image/draw"
	"strings"

	"github.com/nfnt/resize"
)

func MakeMeme(
	img image.Image,
	firstLine, secondLine string,
	watermark bool,
) (image.Image, error) {
	secondLine = strings.ReplaceAll(secondLine, "\n", " ")

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
	var watermarkImg image.Image
	if watermark {
		watermarkImg, err = DrawText("@durich_bot", &TextOptions{
			Color: color.Alpha{A: 100},
			Size:  16,
		})
		if err != nil {
			return nil, err
		}
	}

	imgWidth, imgHeight := img.Bounds().Dx(), img.Bounds().Dy()
	scale := float64(600) / float64(imgWidth)
	imgWidth, imgHeight = int(
		float64(imgWidth)*scale,
	), int(
		float64(imgHeight)*scale,
	)
	if imgWidth < 300 {
		imgWidth = 300
	}

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
	if watermark {
		x0 = drawRect.Max.X - watermarkImg.Bounds().Dx()
		draw.Draw(meme, image.Rect(
			x0,
			drawRect.Min.Y-16-16,
			x0+watermarkImg.Bounds().Dx(),
			drawRect.Min.Y-16-16+watermarkImg.Bounds().Dy(),
		), watermarkImg, image.ZP, draw.Over)
	}

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
