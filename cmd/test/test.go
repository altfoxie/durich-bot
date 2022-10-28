package main

import (
	"fmt"
	"image/png"
	"os"

	"github.com/altfoxie/durich-bot/idraw"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: <input file> <text>")
		os.Exit(1)
	}

	file, err := os.OpenFile(os.Args[1], os.O_RDONLY, 0)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	img, err := png.Decode(file)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	meme, err := idraw.MakeMeme(img, os.Args[2], "", true)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	out, err := os.Create("out.png")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	if err = png.Encode(out, meme); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
