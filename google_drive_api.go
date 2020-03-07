package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
)

type GoogleDrive struct {
	token  string
	client *http.Client
}

type GoogleFileResource struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	WebContentLink string `json:"webContentLink"`
	WebViewLink    string `json:"webViewLink"`
}

type GoogleFileListResource struct {
	Files []GoogleFileResource `json:"files"`
}

type GooglePermissionResource struct {
	Id           string `json:"id"`
	Type         string `json:"type"`
	Role         string `json:"role"`
	EmailAddress string `json:"emailAddress"`
}

func (self *GoogleDrive) addHeader(req *http.Request) {
	req.Header.Add("Authorization", "Bearer "+self.token)
}

func (self *GoogleDrive) do(req *http.Request, result interface{}) error {
	self.addHeader(req)
	res, err := self.client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return errors.New(string(body))
	}
	decoder := json.NewDecoder(res.Body)
	decoder.Decode(result)
	return nil
}

func (self *GoogleDrive) post(req *http.Request, result interface{}) error {
	req.Header.Add("Content-Type", "application/json")
	return self.do(req, result)
}

func (self *GoogleDrive) UploadFilePath(path string) (*GoogleFileResource, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	return self.UploadFileReaderSimple(file)
}

func (self *GoogleDrive) UploadFileReaderMultipart(file io.Reader, args map[string]interface{}) (*GoogleFileResource, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part1, err := writer.CreatePart(textproto.MIMEHeader{
		"Content-Type": []string{"application/json"},
	})

	encoder := json.NewEncoder(part1)
	err = encoder.Encode(args)
	if err != nil {
		return nil, err
	}

	part2, err := writer.CreatePart(textproto.MIMEHeader{
		"Content-Type": []string{"*/*"},
	})
	if err != nil {
		return nil, err
	}

	io.Copy(part2, file)

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		"POST",
		"https://www.googleapis.com/upload/drive/v3/files?uploadType=multipart",
		body,
	)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())

	result := new(GoogleFileResource)

	return result, self.do(req, result)
}

func (self *GoogleDrive) UploadFileReaderSimple(file io.Reader) (*GoogleFileResource, error) {
	req, err := http.NewRequest("POST", "https://www.googleapis.com/upload/drive/v3/files", file)

	if err != nil {
		return nil, err
	}

	result := new(GoogleFileResource)

	return result, self.do(req, result)
}

func (self *GoogleDrive) ListFiles() (*GoogleFileListResource, error) {
	result := new(GoogleFileListResource)
	req, err := http.NewRequest(
		"GET",
		"https://www.googleapis.com/drive/v3/files?fields=kind,files(id,name,webContentLink,webViewLink)",
		nil,
	)

	if err != nil {
		return nil, err
	}

	return result, self.do(req, result)
}

func (self *GoogleDrive) GetFile(id string) (*GoogleFileResource, error) {
	result := new(GoogleFileResource)
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf(
			"https://www.googleapis.com/drive/v3/files/%s?fields=id,name,webContentLink,webViewLink",
			id,
		),
		nil,
	)

	if err != nil {
		return nil, err
	}

	return result, self.do(req, result)
}

func (self *GoogleDrive) CreatePermission(id string, r map[string]interface{}) (*GooglePermissionResource, error) {
	body, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("https://www.googleapis.com/drive/v3/files/%s/permissions?transferOwnership=true", id),
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, err
	}
	result := new(GooglePermissionResource)
	return result, self.post(req, result)
}

func (self *GoogleDrive) DeleteFile(id string) error {
	req, err := http.NewRequest(
		"DELETE",
		fmt.Sprintf("https://www.googleapis.com/drive/v3/files/%s", id),
		nil,
	)
	if err != nil {
		return err
	}
	m := map[string]interface{}{}
	return self.do(req, m)
}
