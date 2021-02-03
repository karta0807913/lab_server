package route

import (
	"crypto/sha256"
	"encoding/base64"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/karta0807913/go_server_utils/serverutil"
	cuserr "github.com/karta0807913/lab_server/error"
	"github.com/karta0807913/lab_server/model"
	"gorm.io/gorm"
)

type APIRouteConfig struct {
	route *gin.RouterGroup
	db    *gorm.DB
}

func bindBody(c *gin.Context, body interface{}) error {
	err := c.ShouldBind(body)
	if err != nil {
		return &cuserr.UserInputError{
			ErrMsg: "body field missing",
		}
	}
	return err
}

func saltPassword(str string) string {
	encoder := sha256.New()
	password := base64.StdEncoding.EncodeToString(
		encoder.Sum(
			[]byte(str),
		),
	)
	return password
}

func APIRouteRegisterHandler(config APIRouteConfig) error {
	route := config.route
	db := config.db

	route.POST("/login", func(c *gin.Context) {
		type Body struct {
			Account  string `json:"account" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		body := new(Body)
		err := bindBody(c, body)
		if err != nil {
			cuserr.GinErrorHandle(err, c)
			return
		}
		password := saltPassword(body.Password)
		log.Println(password)
		var userData model.UserData
		tx := db.Select("ID", "Nickname", "IsAdmin").First(
			&userData,
			"account = ? and password = ?",
			body.Account,
			password,
		)
		if tx.RowsAffected == 0 {
			cuserr.GinErrorHandle(new(cuserr.AccountOrPasswordError), c)
			return
		}
		session := c.MustGet("session").(serverutil.Session)
		session.Set("mem_id", userData.ID)
		session.Set("is_admin", userData.IsAdmin)
		c.JSON(200, userData)
	})

	route.GET("/user", func(c *gin.Context) {
		id, ok := c.GetQuery("user_id")
		if !ok {
			cuserr.GinErrorHandle(&cuserr.UserInputError{
				ErrMsg: "id not found",
			}, c)
			return
		}
		var user model.UserData
		tx := db.Select("ID", "Nickname").Where("id = ?", id).First(&user)
		if tx.RowsAffected == 0 {
			cuserr.GinErrorHandle(&cuserr.UserInputError{
				ErrMsg: "user not found",
			}, c)
			return
		}
		if tx.Error != nil {
			cuserr.GinErrorHandle(tx.Error, c)
			return
		}
		c.JSON(200, user)
	})

	route.GET("/homepage", func(c *gin.Context) {
		var data model.BlogData
		data.Deleted = 0
		db.Select("id", "title", "OwnerID", "Context").Preload("TagList").Where("deleted=? and id=?", data.Deleted, 1).First(&data)
		c.JSON(200, data)
	})

	return nil
}
