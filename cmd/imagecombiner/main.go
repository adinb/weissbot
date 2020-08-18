package main

import (
	"fmt"
	"image/jpeg"
	"os"

	"github.com/adinb/weissbot/pkg/mtg"
)

func main() {
	args := os.Args

	if len(args) != 3 {
		fmt.Println("Usage: imagecombiner <image_1.jpg> <image_1.jpg>")
		os.Exit(1)
	}

	imgPath1 := args[1]
	imgPath2 := args[2]

	imgFile1, err := os.Open(imgPath1)
	if err != nil {
		panic(err.Error())
	}

	imgFile2, err := os.Open(imgPath2)
	if err != nil {
		panic(err.Error())
	}

	img1, err := jpeg.Decode(imgFile1)
	if err != nil {
		panic(err.Error())
	}

	img2, err := jpeg.Decode(imgFile2)
	if err != nil {
		panic(err.Error())
	}

	tiledImage := mtg.TileImagesHorizontally(img1, img2)

	newFile, err := os.Create("result.jpg")
	if err != nil {
		panic(err.Error())
	}

	jpeg.Encode(newFile, tiledImage, nil)
}
