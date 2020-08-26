package route

import (
	"crypto/sha256"
	"encoding/base64"
	"log"

	"github.com/gin-gonic/gin"
	cuserr "github.com/karta0807913/lab_server/error"
	"github.com/karta0807913/lab_server/model"
	"github.com/karta0807913/lab_server/server"
	"gorm.io/gorm"
)

type ApiRouteConfig struct {
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

func ApiRouteRegistHandler(config ApiRouteConfig) error {
	route := config.route
	db := config.db

	route.POST("/login", func(c *gin.Context) {
		log.Println("login")
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
		tx := db.Select("id").First(
			&userData,
			"account = ? and password = ?",
			body.Account,
			password,
		)
		if tx.RowsAffected == 0 {
			cuserr.GinErrorHandle(new(cuserr.AccountOrPasswordError), c)
			return
		}
		session := c.MustGet("session").(server.Session)
		session.Set("mem_id", userData.ID)
		c.JSON(200, map[string]interface{}{
			"state":   0,
			"message": "login success",
		})
	})

	route.POST("/sign_up", func(c *gin.Context) {
		type Body struct {
			Name     string `json:"name" binding:"required"`
			Account  string `json:"account" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		body := new(Body)
		err := bindBody(c, body)
		if err != nil {
			cuserr.GinErrorHandle(err, c)
			return
		}
		tx := db.Select("Nickname", "Account", "Password").Create(&model.UserData{
			Nickname: body.Name,
			Account:  body.Account,
			Password: saltPassword(body.Password),
		})
		if tx.Error != nil {
			log.Println(tx.Error.Error())
			cuserr.GinErrorHandle(new(cuserr.AccountUsed), c)
			return
		}
		c.JSON(200, map[string]interface{}{
			"status": "success",
		})
	})

	route.GET("/me", func(c *gin.Context) {
	})

	return nil
}
