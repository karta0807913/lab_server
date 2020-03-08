package main

import (
	"crypto/sha256"
	"encoding/base64"
	"net/http"
)

func ApiRouteRegistHandler(server *HttpServer, route *SessionServMux) error {
	json_middle := MiddlewareCheckBuilder(new(JsonBodyParser), new(BodyCheck))

	login_sql, err := server.db.Prepare("select mem_id from mem_data where account=? and password=? limit 1")
	if err != nil {
		return err
	}
	route.HandleFunc("/login", func(_w http.ResponseWriter, r *http.Request) {
		type Body struct {
			Account  *string `json:"account"`
			Password *string `json:"password"`
		}
		w := _w.(*ResponseWriter)
		body := new(Body)
		err := json_middle(r, body)

		if err != nil {
			HttpErrorHandle(err, w, r)
			return
		}
		encoder := sha256.New()
		password := base64.StdEncoding.EncodeToString(
			encoder.Sum(
				[]byte(*body.Password),
			),
		)
		result, err := login_sql.Query(body.Account, password)
		if err != nil {
			HttpErrorHandle(err, w, r)
			return
		}
		if !result.Next() {
			HttpErrorHandle(new(AccountOrPasswordError), w, r)
			return
		}
		var mem_id int
		result.Scan(&mem_id)
		w.Session.Set("mem_id", mem_id)
		w.WriteJson(map[string]interface{}{
			"state":   0,
			"message": "login success",
		})
	})

	return nil
}
