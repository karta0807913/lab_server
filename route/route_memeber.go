package route

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/karta0807913/go_server_utils/serverutil"
	cuserr "github.com/karta0807913/lab_server/error"
	"github.com/karta0807913/lab_server/model"
	"gorm.io/gorm"
)

type MemberRouteConfig struct {
	route *gin.RouterGroup
	db    *gorm.DB
}

func MemberRouteRegisterHandler(config MemberRouteConfig) {
	route := config.route
	db := config.db

	route.GET("/me", func(c *gin.Context) {
		session := c.MustGet("session").(serverutil.Session)
		memID := session.Get("mem_id")
		member := model.UserData{}
		tx := db.Select(
			"ID", "Nickname", "IsAdmin",
		).Where("id = ?", memID).First(&member)
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

	route.GET("/logout", func(c *gin.Context) {
		session := c.MustGet("session").(serverutil.Session)
		session.Del("mem_id")
		memID := session.Get("mem_id")
		log.Println(memID)
		c.JSON(200, map[string]interface{}{
			"state": "success",
		})
	})
}
