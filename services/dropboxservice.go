package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/BabyBoChen/bbljfooddiary/utils"
)

type DropboxClient struct {
	accessToken string
}

func NewDropboxClient() (*DropboxClient, error) {
	var dropbox DropboxClient
	err := dropbox.getAccessToken()
	//fmt.Println(dropbox.accessToken)
	return &dropbox, err
}

func (dropbox *DropboxClient) getAccessToken() error {
	endpoint := "https://api.dropboxapi.com/oauth2/token"
	payload := url.Values{}
	envVars := ReadEnvironmentVariables()
	payload.Set("grant_type", "refresh_token")
	payload.Set("refresh_token", envVars.RefreshToken)
	payload.Set("client_id", envVars.DropboxAppKey)
	payload.Set("client_secret", envVars.DropboxAppSecret)
	resp, err := http.Post(endpoint, "application/x-www-form-urlencoded", strings.NewReader(payload.Encode()))

	var resBody []byte
	if err == nil {
		defer resp.Body.Close()
		resBody, err = io.ReadAll(resp.Body)
	}
	var token map[string]interface{}
	if err == nil {
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

// return map[
//
//	client_modified:2023-11-05T12:42:12Z
//	content_hash:91004ff822ee409042e75af0d63e7492d07899b66ccafa0eacc50d2aab3ade41
//	id:id:jRRJCZkdZEIAAAAAAAAADA
//	is_downloadable:true
//	name:1699016250179.jpg
//	path_display:/test123/1699016250179.jpg
//	path_lower:/test123/1699016250179.jpg
//	rev:0160967111472f90000000112aad961
//	server_modified:2023-11-05T12:42:13Z
//	size:233769
//
// ]
func (dropbox *DropboxClient) UploadFile(folderPath string, fileName string, f *os.File) (map[string]interface{}, error) {
	var respJson map[string]interface{}

	urlPath := "https://content.dropboxapi.com/2/files/upload"

	req, err := http.NewRequest("POST", urlPath, f)

	var argsJson []byte
	if err == nil {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", dropbox.accessToken))
		req.Header.Set("Content-Type", "application/octet-stream")
		args := make(map[string]interface{})
		args["path"] = folderPath + "/" + fileName
		mode := make(map[string]interface{})
		mode[".tag"] = "overwrite"
		args["mode"] = mode
		argsJson, err = json.Marshal(args)
	}

	var resp *http.Response
	if err == nil {
		req.Header.Set("Dropbox-API-Arg", string(argsJson))
		client := &http.Client{
			Timeout: time.Second * 30,
		}
		resp, err = client.Do(req)
	}

	var respBody []byte
	if err == nil {
		if resp.StatusCode == 200 {
			defer resp.Body.Close()
			respBody, err = io.ReadAll(resp.Body)
		} else {
			err = errors.New("failed")
		}
	}

	if err == nil {
		err = json.Unmarshal(respBody, &respJson)
	}

	return respJson, err
}

// return map[
//
//	.tag: file (or folder)
//	id:id:a4ayc_80_OEAAAAAAAAAXw
//	name:Prime_Numbers.txt
//	url:https://www.dropbox.com/s/2sn712vy1ovegw8/Prime_Numbers.txt?dl=0
//	...
//
// ]
func (dropbox *DropboxClient) CreateSharedLink(remoteFullPath string) (map[string]interface{}, error) {
	var respJson map[string]interface{}
	var err error

	urlPath := "https://api.dropboxapi.com/2/sharing/create_shared_link_with_settings"

	var body []byte
	payload := make(map[string]interface{})
	payload["path"] = remoteFullPath
	body, err = json.Marshal(payload)

	var req *http.Request
	if err == nil {
		req, err = http.NewRequest("POST", urlPath, bytes.NewBuffer(body))
	}

	var resp *http.Response
	if err == nil {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", dropbox.accessToken))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err = client.Do(req)
	}

	var respBody []byte
	if err == nil {
		if resp.StatusCode == 200 {
			defer resp.Body.Close()
			respBody, err = io.ReadAll(resp.Body)
		} else {
			err = errors.New("failed")
		}
	}

	if err == nil {
		err = json.Unmarshal(respBody, &respJson)
	}

	return respJson, err
}

func (dropbox *DropboxClient) ListFolder(remoteFolderFullPath string) (map[string]interface{}, error) {
	var respJson map[string]interface{}
	var err error

	urlPath := "https://api.dropboxapi.com/2/files/list_folder"

	var body []byte
	payload := make(map[string]interface{})
	payload["include_deleted"] = false
	payload["include_has_explicit_shared_members"] = false
	payload["include_media_info"] = false
	payload["include_mounted_folders"] = false
	payload["include_non_downloadable_files"] = false
	payload["path"] = remoteFolderFullPath
	payload["recursive"] = false
	body, err = json.Marshal(payload)

	var req *http.Request
	if err == nil {
		req, err = http.NewRequest("POST", urlPath, bytes.NewBuffer(body))
	}

	var resp *http.Response
	if err == nil {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", dropbox.accessToken))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err = client.Do(req)
	}

	var respBody []byte
	if err == nil {
		if resp.StatusCode == 200 {
			defer resp.Body.Close()
			respBody, err = io.ReadAll(resp.Body)
		} else {
			err = errors.New("failed")
		}
	}

	if err == nil {
		err = json.Unmarshal(respBody, &respJson)
	}

	var entries []interface{}
	if err == nil {
		entries, err = utils.InterfaceSlice(respJson["entries"])
	}

	if err == nil {
		hasMore := respJson["has_more"].(bool)
		for hasMore && err == nil {
			var more map[string]interface{}
			more, err = dropbox.listFolderMore(respJson["cursor"].(string))

			var moreEntries []interface{}
			if err == nil {
				hasMore = more["has_more"].(bool)
				moreEntries, err = utils.InterfaceSlice(more["entries"])
			}
			if err == nil {
				entries = append(entries, moreEntries...)
			}
		}
	}

	if err == nil {
		respJson["entries"] = entries
	}

	return respJson, err
}

func (dropbox *DropboxClient) listFolderMore(cursor string) (map[string]interface{}, error) {
	var respJson map[string]interface{}
	var err error

	urlPath := "https://api.dropboxapi.com/2/files/list_folder/continue"

	var body []byte
	payload := make(map[string]interface{})
	payload["cursor"] = cursor
	body, err = json.Marshal(payload)
	var req *http.Request
	if err == nil {
		req, err = http.NewRequest("POST", urlPath, bytes.NewBuffer(body))
	}

	var resp *http.Response
	if err == nil {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", dropbox.accessToken))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err = client.Do(req)
	}

	var respBody []byte
	if err == nil {
		if resp.StatusCode == 200 {
			defer resp.Body.Close()
			respBody, err = io.ReadAll(resp.Body)
		} else {
			err = errors.New("failed")
		}
	}

	if err == nil {
		err = json.Unmarshal(respBody, &respJson)
	}

	return respJson, err
}

func (dropbox *DropboxClient) GetSharedLink(remoteFullPath string) ([]string, error) {
	var sharedLink []string
	var err error

	urlPath := "https://api.dropboxapi.com/2/sharing/list_shared_links"

	var body []byte
	payload := make(map[string]interface{})
	payload["path"] = remoteFullPath
	body, err = json.Marshal(payload)

	var req *http.Request
	if err == nil {
		req, err = http.NewRequest("POST", urlPath, bytes.NewBuffer(body))
	}

	var resp *http.Response
	if err == nil {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", dropbox.accessToken))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err = client.Do(req)
	}

	var respBody []byte
	if err == nil {
		if resp.StatusCode == 200 {
			defer resp.Body.Close()
			respBody, err = io.ReadAll(resp.Body)
		} else {
			err = errors.New("failed")
		}
	}

	var respJson map[string]interface{}
	if err == nil {
		err = json.Unmarshal(respBody, &respJson)
	}

	var links []interface{}
	if err == nil {
		links, err = utils.InterfaceSlice(respJson["links"])
	}

	if err == nil {
		sharedLink = make([]string, len(links))
		for i, linkMap := range links {
			link := linkMap.(map[string]interface{})
			sharedLink[i] = link["url"].(string)
		}
	}

	return sharedLink, err
}
