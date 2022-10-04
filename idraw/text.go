package idraw

import (
	"errors"
	"image"
	"image/color"
	"math"

	"github.com/WelcomerTeam/WelcomerImages/pkg/multiface"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

type TextOptions struct {
	Fonts    []*sfnt.Font
	Color    color.Color
	Size     float64
	MaxWidth int
}

func defaultTextOptions(base *TextOptions) *TextOptions {
	opts := &TextOptions{
		Fonts: []*sfnt.Font{defaultFont, defaultEmoji},
		Color: color.White,
		Size:  12,
	}
	if base != nil {
		if len(base.Fonts) != 0 {
			opts.Fonts = base.Fonts
		}
		if base.Color != nil {
			opts.Color = base.Color
		}
		if base.Size != 0 {
			opts.Size = base.Size
		}
		if base.MaxWidth != 0 {
			opts.MaxWidth = base.MaxWidth
		}
	}
	return opts
}

func DrawText(text string, opts *TextOptions) (image.Image, error) {
	opts = defaultTextOptions(opts)

	for size := opts.Size; size > 0; size-- {
		mf := new(multiface.Face)
		for _, f := range opts.Fonts {
			face, err := opentype.NewFace(f, &opentype.FaceOptions{
				Size:    size,
				DPI:     72,
				Hinting: font.HintingFull,
			})
			if err != nil {
				return nil, err
			}
			mf.AddTruetypeFace(face, f)
		}

		width := font.MeasureString(mf, text).Ceil()
		if opts.MaxWidth != 0 && width > opts.MaxWidth {
			size = math.Ceil(float64(opts.MaxWidth) / float64(width) * size)
			continue
		}
		height := mf.Metrics().Ascent.Ceil() + mf.Metrics().Descent.Ceil()

		drawer := &font.Drawer{
			Dst: image.NewRGBA(
				image.Rect(0, 0, width, height),
			),
			Src:  image.NewUniform(opts.Color),
			Face: mf,
			Dot:  fixed.P(0, mf.Metrics().Ascent.Ceil()),
		}
		drawer.DrawString(text)
		return drawer.Dst, nil
	}

	return nil, errors.New("failed to choose a font size")
}
