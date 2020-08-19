package cuserr

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

type IsNotJsonError struct {
	error
}

func (IsNotJsonError) Error() string {
	return "not a json body"
}

type UserInputError struct {
	error
	ErrMsg string
}

func (self UserInputError) Error() string {
	return self.ErrMsg
}

type FileUploadError struct {
	error
	ErrMsg string
}

func (self FileUploadError) Error() string {
	return self.ErrMsg
}

// TODO: add log and more error response
func HttpErrorHandle(err error, w http.ResponseWriter, r *http.Request) {
	switch err := err.(type) {
	default:
		w.WriteHeader(500)
		w.Write([]byte("server error"))
		log.Println(err.Error())
		break
	case *PleasLoginError,
		*AccountOrPasswordError,
		*IsNotJsonError,
		*UserInputError:
		w.WriteHeader(403)
		w.Write([]byte(err.Error()))
		break
	case *FileUploadError:
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}
}
