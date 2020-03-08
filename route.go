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

	file_route := LoginChecker{
		Handler: http.StripPrefix("/file", http.NewServeMux()),
	}
	router.Handle("/file/", file_route)
	FileRouteRegistHandler(server, file_route.Handler.(*http.ServeMux))

	api_route := http.NewServeMux()
	ApiRouteRegistHandler(server, api_route)
}
