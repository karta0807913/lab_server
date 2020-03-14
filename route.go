package main

import "fmt"
import "net/http"

type LoginChecker struct {
	http.Handler
}

func (self LoginChecker) ServeHTTP(_w http.ResponseWriter, r *http.Request) {
	w := _w.(*ResponseWriter)
	fmt.Println(w.Session.session)
	if w.Session.Get("mem_id") == nil {
		HttpErrorHandle(new(PleasLoginError), w, r)
		return
	}
	self.Handler.ServeHTTP(w, r)
}

func route(server *HttpServer) {
	router := server.Handler.(*SessionServMux)
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
	FileRouteRegistHandler(server, file_route)

	api_route := http.NewServeMux()
	ApiRouteRegistHandler(server, api_route)
	router.Handle("/api/", http.StripPrefix("/api", api_route))
}
