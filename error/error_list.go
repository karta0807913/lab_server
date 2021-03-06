package cuserr

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(404, gin.H{
				"message": "record not found",
			})
		} else if errors.Is(err, io.EOF) {
			c.AbortWithStatusJSON(403, gin.H{
				"message": "http body required",
			})
		} else {
			c.AbortWithStatusJSON(500, gin.H{
				"message": "unknow error",
			})
			log.Println(err.Error())
		}
	case *PleasLoginError,
		*AccountOrPasswordError,
		*AccountUsed,
		*IsNotJsonError,
		*UserInputError:
		c.AbortWithStatusJSON(403, gin.H{
			"message": err.Error(),
		})
	case *FileUploadError:
		c.AbortWithStatusJSON(500, gin.H{
			"message": err.Error(),
		})
	}
}
