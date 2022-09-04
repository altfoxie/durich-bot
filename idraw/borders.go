package idraw

import (
	"image"
	"image/draw"
)

type BordersOptions struct {
	Color *image.Uniform
	Size  int
}

func defaultBordersOptions(base *BordersOptions) *BordersOptions {
	opts := &BordersOptions{
		Color: image.White,
		Size:  5,
	}
	if base != nil {
		if base.Color != nil {
			opts.Color = base.Color
		}
		if base.Size != 0 {
			opts.Size = base.Size
		}
	}
	return opts
}

func DrawBorders(img draw.Image, rect image.Rectangle, opts *BordersOptions) {
	opts = defaultBordersOptions(opts)
	borders := []image.Rectangle{
		image.Rect(rect.Min.X, rect.Min.Y-opts.Size, rect.Max.X, rect.Min.Y),
		image.Rect(rect.Max.X, rect.Max.Y+opts.Size, rect.Min.X, rect.Max.Y),
		image.Rect(
			rect.Min.X-opts.Size,
			rect.Min.Y-opts.Size,
			rect.Min.X,
			rect.Max.Y+opts.Size,
		),
		image.Rect(
			rect.Max.X+opts.Size,
			rect.Min.Y-opts.Size,
			rect.Max.X,
			rect.Max.Y+opts.Size,
		),
	}
	for _, rect := range borders {
		draw.Draw(img, rect, opts.Color, image.ZP, draw.Over)
	}
}
