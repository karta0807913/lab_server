package route

import (
	"net/http"

	cuserr "github.com/karta0807913/lab_server/error"
	"github.com/karta0807913/lab_server/server"
)

type LoginChecker struct {
	http.Handler
}

func (self LoginChecker) ServeHTTP(_w http.ResponseWriter, r *http.Request) {
	w := _w.(*server.ResponseWriter)
	// fmt.Println(w.Session.session)
	if w.Session.Get("mem_id") == nil {
		cuserr.HttpErrorHandle(new(cuserr.PleasLoginError), w, r)
		return
	}
	self.Handler.ServeHTTP(w, r)
}

func Route(serv *server.HttpServer, upload_path string) {
	router := serv.Handler.(*server.SessionServMux)
	// router.HandleFunc("/", func(_w http.ResponseWriter, r *http.Request) {
	// 	w := _w.(*ResponseWriter)
	//     fmt.Println(w.Session.Get("A"))
	//     if w.Session.Get("A") != nil && w.Session.Get("A").(string) == "B" {
	//         w.Session.Set("A", "C")
	//     } else {
	//         w.Session.Set("A", "B")
	//     }
	// 	w.Write([]byte("Hello"))
	// })

	file_route := http.NewServeMux()
	file_login_check := LoginChecker{
		Handler: http.StripPrefix("/file", file_route),
	}
	router.Handle("/file/", file_login_check)
	FileRouteRegistHandler(serv, file_route, upload_path)

	api_route := http.NewServeMux()
	ApiRouteRegistHandler(serv, api_route)
	router.Handle("/api/", http.StripPrefix("/api", api_route))
}
