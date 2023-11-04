package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/BabyBoChen/bbljfooddiary/secrets"
)

const dropboxApi string = "https://api.dropboxapi.com/oauth2/token"

type DropboxClient struct {
	accessToken string
}

func NewDropboxClient() (DropboxClient, error) {
	var dropbox DropboxClient
	err := dropbox.getAccessToken()
	//fmt.Println(dropbox.accessToken)
	return dropbox, err
}

func (dropbox *DropboxClient) getAccessToken() error {
	payload := url.Values{}
	payload.Set("grant_type", "refresh_token")
	payload.Set("refresh_token", secrets.RefreshToken)
	payload.Set("client_id", secrets.DropboxAppKey)
	payload.Set("client_secret", secrets.DropboxAppSecret)
	resp, err := http.Post(dropboxApi, "application/x-www-form-urlencoded", strings.NewReader(payload.Encode()))
	var token map[string]interface{}
	if err == nil {
		defer resp.Body.Close()
		resBody, _ := io.ReadAll(resp.Body)
		err = json.Unmarshal(resBody, &token)
	}
	if err == nil {
		dropbox.accessToken = token["access_token"].(string)
	}
	if err != nil {
		fmt.Println(err)
	}
	return err
}
