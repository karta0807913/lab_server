package route

import (
	"log"

	"github.com/gin-gonic/gin"
	cuserr "github.com/karta0807913/lab_server/error"
	"github.com/karta0807913/lab_server/model"
	"github.com/karta0807913/lab_server/server"
	"gorm.io/gorm"
)

type MemberRouteConfig struct {
	route *gin.RouterGroup
	db    *gorm.DB
}

func MemberRouteRegistHandler(config MemberRouteConfig) {
	route := config.route
	db := config.db

	route.GET("/me", func(c *gin.Context) {
		session := c.MustGet("session").(server.Session)
		mem_id := session.Get("mem_id")
		member := model.UserData{}
		tx := db.Select("id", "nickname").Where("id = ?", mem_id).First(&member)
		if tx.Error != nil {
			cuserr.GinErrorHandle(tx.Error, c)
			return
		}
		c.JSON(200, member)
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
		tx := db.Select("id", "nickname").Where("id = ?", id).First(&user)
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

	route.GET("/logout", func(c *gin.Context) {
		session := c.MustGet("session").(server.Session)
		session.Del("mem_id")
		mem_id := session.Get("mem_id")
		log.Println(mem_id)
		c.JSON(200, map[string]interface{}{
			"state": "success",
		})
	})
}
