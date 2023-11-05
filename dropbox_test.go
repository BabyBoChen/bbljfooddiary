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
	imgs := make([]string, 6)
	imgs[0] = "./tmp/1699016250179.jpg"
	imgs[1] = "./tmp/1699016250179 (copy).jpg"
	imgs[2] = "./tmp/1699016250179 (another copy).jpg"
	imgs[3] = "./tmp/1699016250179 (3rd copy).jpg"
	imgs[4] = "./tmp/1699016250179 (4th copy).jpg"
	imgs[5] = "./tmp/1699016250179 (5th copy).jpg"
	for _, fn := range imgs {
		f, _ := os.Open(fn)
		resp, _ := dropbox.UploadFile("/testupload", f)
		path := resp["path_lower"].(string)
		dropbox.CreateSharedLink(path)
	}
}

func TestListFolder(t *testing.T) {
	dropbox, _ := services.NewDropboxClient()
	ls, _ := dropbox.ListFolder("/testupload")
	fmt.Println(ls)
}

func TestGetSharedLink(t *testing.T) {
	dropbox, _ := services.NewDropboxClient()
	links, err := dropbox.GetSharedLink("/testupload/1699016250179 (another copy).jpg")
	if err == nil {
		fmt.Println(links)
	} else {
		fmt.Println(err)
	}
}
