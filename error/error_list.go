package cuserr

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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

type FileNotFoundError struct {
	error
	FileId string
}

func (self FileNotFoundError) Error() string {
	return fmt.Sprintf("file %s not found", self.FileId)
}

type AccountUsed struct {
	error
}

func (self AccountUsed) Error() string {
	return "Account used"
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
		*AccountUsed,
		*UserInputError:
		w.WriteHeader(403)
		w.Write([]byte(err.Error()))
		break
	case *FileUploadError:
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}
}

func GinErrorHandle(err error, c *gin.Context) {
	switch err := err.(type) {
	default:
		c.Writer.WriteHeader(500)
		c.Writer.Write([]byte("server error"))
		log.Println(err.Error())
		break
	case *PleasLoginError,
		*AccountOrPasswordError,
		*AccountUsed,
		*IsNotJsonError,
		*UserInputError:
		c.Writer.WriteHeader(403)
		c.Writer.Write([]byte(err.Error()))
		break
	case *FileUploadError:
		c.Writer.WriteHeader(500)
		c.Writer.Write([]byte(err.Error()))
	}
}
