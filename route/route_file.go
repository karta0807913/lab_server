package route

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/karta0807913/go_server_utils/serverutil"
	cuserr "github.com/karta0807913/lab_server/error"
	"github.com/karta0807913/lab_server/model"
	"gorm.io/gorm"
)

const (
	FileWaiting     = "0"
	FileDownloading = "1"
	FileDone        = "2"
)

type FileRouteConfig struct {
	route      *gin.RouterGroup
	uploadPath string
	db         *gorm.DB
}

func FileRouteRegisterHandler(config FileRouteConfig) {
	uploadPath := config.uploadPath
	route := config.route
	db := config.db

	route.POST("/upload", func(c *gin.Context) {
		type UploadParameter struct {
			Filename string `json:"filename"`
			TagID    *uint  `json:"tag_id"`
		}

		reader, err := c.Request.MultipartReader()
		if err != nil {
			cuserr.GinErrorHandle(err, c)
			return
		}
		part, err := reader.NextPart()
		if err == io.EOF {
			cuserr.GinErrorHandle(cuserr.UserInputError{
				ErrMsg: "format error",
			}, c)
			return
		}
		if part.Header.Get("Content-Type") != "application/json" {
			cuserr.GinErrorHandle(cuserr.UserInputError{
				ErrMsg: "first block must be json type",
			}, c)
			return
		}
		decoder := json.NewDecoder(part)
		param := new(UploadParameter)
		err = decoder.Decode(param)
		if err != nil {
			cuserr.GinErrorHandle(err, c)
			return
		}
		if param.TagID == nil {
			cuserr.GinErrorHandle(cuserr.UserInputError{
				ErrMsg: "Tag id missing",
			}, c)
			return
		}

		part, err = reader.NextPart()
		if err == io.EOF {
			cuserr.GinErrorHandle(cuserr.UserInputError{
				ErrMsg: "file missing",
			}, c)
			return
		}

		crypto := sha512.New()
		fileReader := io.TeeReader(part, crypto)

		tmpFile, err := ioutil.TempFile(
			path.Join(uploadPath, "temp"), "uploading_*",
		)
		if err != nil {
			cuserr.GinErrorHandle(cuserr.FileUploadError{
				ErrMsg: "Can't create temporary file",
			}, c)
			return
		}

		io.Copy(tmpFile, fileReader)
		sum := crypto.Sum(nil)
		fileHash := base64.URLEncoding.EncodeToString(sum)

		path := path.Join(uploadPath, fileHash)
		_, err = os.Stat(path)
		if os.IsNotExist(err) {
			os.Rename(tmpFile.Name(), path)
		} else {
			os.Remove(tmpFile.Name())
		}

		fileData := model.FileData{
			Filename:    param.Filename,
			FileHash:    fileHash,
			UserID:      uint(c.MustGet("session").(serverutil.Session).Get("mem_id").(float64)),
			ContextType: part.Header.Values("Content-Type")[0],
		}
		tx := db.Select("file_name", "file_hash", "user_id").Create(&fileData)
		if tx.Error != nil {
			cuserr.GinErrorHandle(tx.Error, c)
			return
		}

		c.JSON(200, map[string]interface{}{
			"msg":     "request accept",
			"file_id": fileData.ID,
		})
	})

	route.GET("/list", func(c *gin.Context) {
		var fileList []model.FileData
		id, ok := c.GetQuery("id")
		var tx *gorm.DB
		if ok {
			tx = db.Select([]string{
				"id", "filename", "context_type", "user_id",
			}).Where("deleted = 0 and id = ?", id).Where(&fileList)
		} else {
			tx = db.Select([]string{
				"id", "filename", "context_type", "user_id",
			}).Where("deleted = 0").Find(&fileList)
		}
		if tx.Error != nil {
			cuserr.GinErrorHandle(tx.Error, c)
			return
		}

		c.JSON(200, fileList)
	})

	route.GET("/download", func(c *gin.Context) {
		fileInfo := model.FileData{}
		id, ok := c.GetQuery("id")
		if !ok {
			cuserr.GinErrorHandle(&cuserr.UserInputError{
				ErrMsg: "File id missing",
			}, c)
			return
		}
		tx := db.Select("Filename", "FileHash").Where("deleted = 0 and id = ?", id).First(&fileInfo)
		if tx.Error != nil {
			cuserr.GinErrorHandle(tx.Error, c)
			return
		}
		filePath := path.Join(uploadPath, fileInfo.FileHash)
		_, err := os.Stat(filePath)
		if err != nil {
			cuserr.GinErrorHandle(&cuserr.FileNotFoundError{
				FileId: id,
			}, c)
			return
		}
		c.Header("Content-Type", fileInfo.ContextType)
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileInfo.Filename))
		c.File(filePath)
	})
}
