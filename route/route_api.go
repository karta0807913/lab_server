package route

import (
	"crypto/sha256"
	"encoding/base64"
	"net/http"

	cuserr "github.com/karta0807913/lab_server/error"
	"github.com/karta0807913/lab_server/model"
	"github.com/karta0807913/lab_server/server"
)

func ApiRouteRegistHandler(serv *server.HttpServer, route *http.ServeMux) error {
	json_middle := server.MiddlewareCheckBuilder(new(server.JsonBodyParser), new(server.BodyCheck))

	route.HandleFunc("/login", func(_w http.ResponseWriter, r *http.Request) {
		type Body struct {
			Account  *string `json:"account"`
			Password *string `json:"password"`
		}
		w := _w.(*server.ResponseWriter)
		body := new(Body)
		err := json_middle(r, body)

		if err != nil {
			cuserr.HttpErrorHandle(err, w, r)
			return
		}
		encoder := sha256.New()
		password := base64.StdEncoding.EncodeToString(
			encoder.Sum(
				[]byte(*body.Password),
			),
		)
		println(password)
		var userData model.UserData
		tx := serv.DB().Select("id").First(&userData, "account = ? and password = ?", body.Account, password)
		if tx.Error != nil {
			cuserr.HttpErrorHandle(tx.Error, w, r)
			return
		}
		if tx.RowsAffected == 0 {
			cuserr.HttpErrorHandle(new(cuserr.AccountOrPasswordError), w, r)
			return
		}
		w.Session.Set("mem_id", userData.ID)
		w.WriteJson(map[string]interface{}{
			"state":   0,
			"message": "login success",
		})
	})

	return nil
}
