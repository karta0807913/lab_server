package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

const (
	FileWaiting     = "0"
	FileDownloading = "1"
	FileDone        = "2"
)

func FileRouteRegistHandler(server *HttpServer, route *http.ServeMux) {
	type UploadParameter struct {
		Filename string `json:"filename"`
	}

	type UploadFileRequest struct {
		file io.Reader
		info map[string]interface{}
		id   int32
	}
	type UploadFileResult struct {
		err      error
		id       int32
		resource *GoogleFileResource
	}

	upload_stream := make(chan UploadFileRequest)
	result_stream := make(chan UploadFileResult)

	UploadStream := func(request_stream chan UploadFileRequest, result_stream chan UploadFileResult) {
		for request, ok := <-request_stream; ok; request, ok = <-request_stream {
			_, err := server.db.Exec(
				"update file_data set `state`=? where id=?",
				FileDownloading, request.id,
			)
			if err != nil {
				log.Println(err)
				_, err := server.db.Exec(
					"delete from file_data where id=?",
					request.id,
				)
				if err != nil {
					log.Println(err)
				}
				continue
			}
			info, err := server.drive.UploadFileReaderMultipart(request.file, request.info)
			result_stream <- UploadFileResult{
				err:      err,
				resource: info,
				id:       request.id,
			}
		}
	}

	ProcessResult := func(result_stream chan UploadFileResult) {
		for response, ok := <-result_stream; ok; response, ok = <-result_stream {
			if response.err != nil {
				resource := response.resource
				_, err := server.db.Exec(
					"update file_data set `google_id`=?, `state`=? where id=?",
					resource.Id, FileDone, response.id,
				)
				if err != nil {
					log.Println(err)
					result_stream <- response
					continue
				}
				_, err = server.drive.CreatePermission(resource.Id, map[string]interface{}{
					"type": "anyone",
					"role": "writer",
				})
				if err != nil {
					log.Println(err)
					result_stream <- response
					continue
				}
			} else {
				log.Println(response.err)
			}
		}
	}

	go UploadStream(upload_stream, result_stream)
	go ProcessResult(result_stream)

	checkPostMethod := MethodCheck{
		method: "POST",
	}

	route.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		err := checkPostMethod.Handle(r, nil)
		if err != nil {
			HttpErrorHandle(err, w, r)
			return
		}
		reader, err := r.MultipartReader()
		if err != nil {
			HttpErrorHandle(err, w, r)
			return
		}
		part, err := reader.NextPart()
		if err == io.EOF {
			HttpErrorHandle(UserInputError{
				err_msg: "format error",
			}, w, r)
			return
		}
		if part.Header.Get("Content-Type") != "application/json" {
			HttpErrorHandle(UserInputError{
				err_msg: "first block must be json type",
			}, w, r)
			return
		}
		decoder := json.NewDecoder(part)
		param := new(UploadParameter)
		err = decoder.Decode(param)
		if err != nil {
			HttpErrorHandle(err, w, r)
			return
		}

		part, err = reader.NextPart()
		if err == io.EOF {
			HttpErrorHandle(UserInputError{
				err_msg: "file missing",
			}, w, r)
			return
		}

		result, err := server.db.Exec("insert into file_data (`state`) values (`?`)", FileWaiting)
		if err != nil {
			HttpErrorHandle(err, w, r)
			return
		}

		info := map[string]interface{}{
			"filename": param.Filename,
		}
		encoder := json.NewEncoder(w)
		last_id, err := result.LastInsertId()
		if err != nil {
			HttpErrorHandle(err, w, r)
			return
		}

		upload_stream <- UploadFileRequest{
			file: part,
			info: map[string]interface{}{
				"filename": param.Filename,
				"mimeType": info["Content-Type"],
				"parents":  []string{Config.google_file_parent},
			},
		}

		w.WriteHeader(200)
		encoder.Encode(map[string]interface{}{
			"msg":     "request accept",
			"file_id": last_id,
		})
	})

	checkGetMethod := MethodCheck{
		method: "GET",
	}
	route.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		err := checkGetMethod.Handle(r, nil)
		if err != nil {
			HttpErrorHandle(err, w, r)
			return
		}

		info_list, err := server.drive.ListFiles()
		if err != nil {
			HttpErrorHandle(err, w, r)
			return
		}

		w.WriteHeader(200)
		encoder := json.NewEncoder(w)
		err = encoder.Encode(info_list)
		if err != nil {
			HttpErrorHandle(err, w, r)
			return
		}
	})

	route.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		err := checkGetMethod.Handle(r, nil)
		if err != nil {
			HttpErrorHandle(err, w, r)
			return
		}
		values := r.URL.Query()
		id := values.Get("id")
		res, err := server.db.Query("select google_id, state from file_data where id=? limit 1", id)
		if err != nil {
			HttpErrorHandle(err, w, r)
			return
		}
		if !res.Next() {
			w.WriteHeader(403)
			w.Write([]byte(`{ "msg": "id not found or upload failed" }`))
			return
		}
		data, err := res.Columns()
		if err != nil {
			HttpErrorHandle(err, w, r)
			return
		}
		if data[1] != FileDone {
			w.WriteHeader(406)
			w.Write([]byte(`{ "msg": "file uploading" }`))
			return
		}
		file_info, err := server.drive.GetFile(data[0])
		if err != nil {
			HttpErrorHandle(err, w, r)
			return
		}

		encoder := json.NewEncoder(w)
		err = encoder.Encode(file_info)
		if err != nil {
			HttpErrorHandle(err, w, r)
			return
		}
	})
}
