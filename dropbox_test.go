package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/BabyBoChen/bbljfooddiary/services"
)

func TestDropBox(t *testing.T) {
	services.NewDropboxClient()
}

func TestUpload(t *testing.T) {
	dropbox, _ := services.NewDropboxClient()
	f, _ := os.Open("./tmp/1699016250179.jpg")
	resp, _ := dropbox.UploadFile("/test123", f)
	fmt.Println(resp)
	path := resp["path_lower"].(string)
	resp, _ = dropbox.CreateSharedLink(path)
	fmt.Println(resp)
}
