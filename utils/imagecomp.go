package utils

import (
	"fmt"
	"image"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
)

func ResizeImage(img image.Image, ext string, maxWidth int) string {
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()
	var comp *image.NRGBA
	if w > h {
		if maxWidth > w {
			maxWidth = w
		}
		comp = imaging.Resize(img, maxWidth, 0, imaging.CatmullRom)
	} else {
		if maxWidth > h {
			maxWidth = h
		}
		comp = imaging.Resize(img, 0, maxWidth, imaging.CatmullRom)
	}
	uid := uuid.New()
	if len(ext) > 0 && ext[0] != '.' {
		ext = "." + ext
	}
	newFileName := fmt.Sprintf("./tmp/%s%s", uid.String(), ext)
	imaging.Save(comp, newFileName)
	return newFileName
}
