package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
)

func TestUpload(t *testing.T) {
	token, err := google_token("./e539-lab-web-dd38239bcca2.json",
		"https://www.googleapis.com/auth/drive")
	if err != nil {
		log.Fatal(err)
	}

	drive := GoogleDrive{token, &http.Client{}}
	file, err := drive.UploadFilePath("./main.go")
	if err != nil {
		log.Fatal(err)
	}
	_, err = drive.CreatePermission(file.Id, map[string]interface{}{
		"type": "anyone",
		"role": "reader",
	})
	if err != nil {
		log.Fatal(err)
	}
	file, err = drive.GetFile(file.Id)
	if err != nil {
		log.Fatal(err)
	}
	println(file.WebViewLink)
	res, err := http.Get(file.WebContentLink)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	result, err := os.OpenFile("./main.go", os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	actual, err := ioutil.ReadAll(result)
	if err != nil {
		log.Fatal(err)
	}
	if string(actual) != string(body) {
		t.Errorf("File content not match, get \n%s", string(actual))
	}
}
