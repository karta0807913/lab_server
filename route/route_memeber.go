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

	route.PUT("/change_password", func(c *gin.Context) {
		session := c.MustGet("session").(serverutil.Session)
		memID := session.Get("mem_id")
		type Body struct {
			OldPassword string `json:"old_password" binding:"required"`
			NewPassword string `json:"new_password" binding:"required"`
		}
		var body Body
		err := c.ShouldBindJSON(&body)
		if err != nil {
			cuserr.GinErrorHandle(err, c)
			return
		}

		var userData model.UserData
		err = db.Where("id=? and password=?", memID, saltPassword(body.OldPassword)).Find(&userData).Error
		if err != nil || userData.ID == 0 {
			log.Println(err)
			c.JSON(403, gin.H{
				"state":   "failed",
				"message": "password not correct",
			})
			return
		}
		userData.Password = saltPassword(body.NewPassword)
		err = db.Model(&userData).Select("Password").Updates(&userData).Error
		if err != nil {
			cuserr.GinErrorHandle(err, c)
			return
		}
		c.JSON(200, gin.H{
			"message": "success",
		})
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
