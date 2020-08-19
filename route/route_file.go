package route

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path"

	cuserr "github.com/karta0807913/lab_server/error"
	"github.com/karta0807913/lab_server/model"
	"github.com/karta0807913/lab_server/server"
)

const (
	FileWaiting     = "0"
	FileDownloading = "1"
	FileDone        = "2"
)

func FileRouteRegistHandler(serv *server.HttpServer, route *http.ServeMux, upload_path string) {
	type UploadParameter struct {
		Filename string `json:"filename"`
	}

	checkPostMethod := server.MethodCheck{
		Method: "POST",
	}

	route.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		err := checkPostMethod.Handle(r, nil)
		if err != nil {
			cuserr.HttpErrorHandle(err, w, r)
			return
		}
		reader, err := r.MultipartReader()
		if err != nil {
			cuserr.HttpErrorHandle(err, w, r)
			return
		}
		part, err := reader.NextPart()
		if err == io.EOF {
			cuserr.HttpErrorHandle(cuserr.UserInputError{
				ErrMsg: "format error",
			}, w, r)
			return
		}
		if part.Header.Get("Content-Type") != "application/json" {
			cuserr.HttpErrorHandle(cuserr.UserInputError{
				ErrMsg: "first block must be json type",
			}, w, r)
			return
		}
		decoder := json.NewDecoder(part)
		param := new(UploadParameter)
		err = decoder.Decode(param)
		if err != nil {
			cuserr.HttpErrorHandle(err, w, r)
			return
		}

		part, err = reader.NextPart()
		if err == io.EOF {
			cuserr.HttpErrorHandle(cuserr.UserInputError{
				ErrMsg: "file missing",
			}, w, r)
			return
		}

		crypto := sha512.New()
		file_reader := io.TeeReader(part, crypto)
		fileHash := base64.URLEncoding.EncodeToString(crypto.Sum(nil))

		path := path.Join(upload_path, fileHash)
		_, err = os.Stat(path)
		if os.IsNotExist(err) {
			file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0640)
			if err != nil {
				cuserr.HttpErrorHandle(cuserr.FileUploadError{
					ErrMsg: "Can't create file",
				}, w, r)
				return
			}
			defer file.Close()
			io.Copy(file, file_reader)
		}

		fileData := model.FileData{
			Filename: param.Filename,
			FileHash: fileHash,
		}
		tx := serv.DB().Create(fileData)
		if tx.Error != nil {
			cuserr.HttpErrorHandle(tx.Error, w, r)
			return
		}

		encoder := json.NewEncoder(w)

		w.WriteHeader(200)
		encoder.Encode(map[string]interface{}{
			"msg":     "request accept",
			"file_id": fileData.ID,
		})
	})

	checkGetMethod := server.MethodCheck{
		Method: "GET",
	}
	route.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		err := checkGetMethod.Handle(r, nil)
		if err != nil {
			cuserr.HttpErrorHandle(err, w, r)
			return
		}

		var file_list []model.FileData
		tx := serv.DB().Select([]string{"id", "filename"}).Where("deleted = 0").Find(&file_list)
		if tx.Error != nil {
			cuserr.HttpErrorHandle(tx.Error, w, r)
			return
		}

		w.WriteHeader(200)
		encoder := json.NewEncoder(w)
		err = encoder.Encode(file_list)
		if err != nil {
			cuserr.HttpErrorHandle(err, w, r)
			return
		}
	})

	route.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		err := checkGetMethod.Handle(r, nil)
		if err != nil {
			cuserr.HttpErrorHandle(err, w, r)
			return
		}
		values := r.URL.Query()
		id := values.Get("id")
		var fileData model.FileData
		tx := serv.DB().First(&fileData, id)
		if tx.Error != nil {
			cuserr.HttpErrorHandle(tx.Error, w, r)
			return
		}
		if tx.RowsAffected == 0 {
			cuserr.HttpErrorHandle(cuserr.UserInputError{
				ErrMsg: "File Not Found",
			}, w, r)
		}
		encoder := json.NewEncoder(w)
		err = encoder.Encode(fileData)
		if err != nil {
			cuserr.HttpErrorHandle(err, w, r)
			return
		}
	})
}
