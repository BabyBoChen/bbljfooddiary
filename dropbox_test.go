package main

import (
	"testing"

	"github.com/BabyBoChen/bbljfooddiary/services"
)

func TestDropBox(t *testing.T) {
	services.NewDropboxClient()
}
