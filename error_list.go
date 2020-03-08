package main

import (
	"log"
	"net/http"
)

type AccountOrPasswordError struct {
	error
}

func (AccountOrPasswordError) Error() string {
	return "AccountOrPasswordError"
}

type PleasLoginError struct {
	error
}

func (PleasLoginError) Error() string {
	return "Please Login"
}

func HttpErrorHandle(err error, w http.ResponseWriter, r *http.Request) {
	switch err := err.(type) {
	default:
		w.WriteHeader(500)
		w.Write([]byte("server error"))
		log.Println(err.Error())
		break
	case PleasLoginError:
	case AccountOrPasswordError:
	case UserInputError:
		w.WriteHeader(403)
		w.Write([]byte(err.Error()))
		break
	}
}
