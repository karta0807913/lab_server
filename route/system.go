package route

import (
	"log"

	"github.com/gin-gonic/gin"
	cuserr "github.com/karta0807913/lab_server/error"
	"github.com/karta0807913/lab_server/model"
)

func APIRouteSystemHandler(config APIRouteConfig) error {
	route := config.route
	db := config.db

	route.POST("/new_account", func(c *gin.Context) {
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
		tx := db.Select(
			"Nickname", "Account", "Password",
		).Create(&model.UserData{
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
	return nil
}
