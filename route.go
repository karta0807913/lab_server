package main

import "net/http"

func route(server *HttpServer) {
	router := server.Handler.(*http.ServeMux)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	})

	file_route := new(http.ServeMux)
	router.Handle("/file", file_route)
	FileRouteRegistHandler(server, file_route)
}
