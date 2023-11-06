package main

import (
	"image"
	"os"
	"testing"

	"github.com/BabyBoChen/bbljfooddiary/utils"
)

func TestImageCo(t *testing.T) {
	f, _ := os.Open("./tmp/1699016250179.jpg")
	img, ext, _ := image.Decode(f)
	utils.ResizeImage(img, ext, 800)
}
