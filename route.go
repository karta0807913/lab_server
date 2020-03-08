package main

import "net/http"

type LoginChecker struct {
	http.Handler
}

func (self LoginChecker) ServeHTTP(_w http.ResponseWriter, r *http.Request) {
	w := _w.(*ResponseWriter)
	if w.Session.Get("mem_id") == nil {
		HttpErrorHandle(new(PleasLoginError), w, r)
		return
	}
	self.Handler.ServeHTTP(w, r)
}

func route(server *HttpServer) {
	router := server.Handler.(*SessionServMux)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	})

	file_route := http.NewServeMux()
	file_login_check := LoginChecker{
		Handler: http.StripPrefix("/file/", file_route),
	}
	router.Handle("/file/", file_login_check)
	FileRouteRegistHandler(server, file_route)

	api_route := http.NewServeMux()
	router.Handle("/api/", http.StripPrefix("/api/", api_route))
	ApiRouteRegistHandler(server, api_route)
}
