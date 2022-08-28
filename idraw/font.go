package idraw

import (
	_ "embed"

	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
)

var (
	// Let's pretend that you don't know this font?
	//go:embed font.ttf
	defaultFontBytes []byte
	defaultFont      *sfnt.Font
)

func init() {
	var err error
	defaultFont, err = opentype.Parse(defaultFontBytes)
	if err != nil {
		panic(err)
	}
}
