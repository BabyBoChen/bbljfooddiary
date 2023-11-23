package main

import (
	"fmt"
	"image"
	"os"
	"testing"

	"github.com/BabyBoChen/bbljfooddiary/utils"
)

func TestImageCo(t *testing.T) {
	f, _ := os.Open("./tmp/1699016250179.jpg")
	img, ext, _ := image.Decode(f)
	saveFileName := fmt.Sprintf("test.%s", ext)
	utils.ResizeImage(img, saveFileName, 800)
}
