package route

import (
	"github.com/gin-gonic/gin"
	"github.com/karta0807913/go_server_utils/serverutil"
	cuserr "github.com/karta0807913/lab_server/error"
	"github.com/karta0807913/lab_server/model"
)

func AdminRouteRegisterHandler(config APIRouteConfig) {
	db := config.db
	route := config.route

	AdminCheck := func(c *gin.Context) {
		session := c.MustGet("session").(serverutil.Session)
		if session.Get("is_admin").(bool) {
			c.Next()
		} else {
			c.AbortWithStatusJSON(403, gin.H{
				"message": "permission denied",
			})
		}
	}

	route.Use(AdminCheck)

	route.GET("/member", func(c *gin.Context) {
		var user model.UserData
		data, err := user.Find(c, db)
		if err != nil {
			cuserr.GinErrorHandle(err, c)
			return
		}
		c.JSON(200, data)
	})

	route.PUT("/member", func(c *gin.Context) {
		var user model.UserData
		err := user.Update(c, db)
		if err != nil {
			cuserr.GinErrorHandle(err, c)
			return
		}
		c.JSON(200, gin.H{
			"message": "success",
		})
	})

	route.GET("/tag", func(c *gin.Context) {
		var tag []model.BlogTag
		db.Find(&tag)
		c.JSON(200, tag)
	})

	route.POST("/tag", func(c *gin.Context) {
		var tag model.BlogTag
		err := tag.Create(c, db)
		if err != nil {
			cuserr.GinErrorHandle(err, c)
			return
		}
		c.JSON(200, tag)
	})

	route.PUT("/tag", func(c *gin.Context) {
		var tag model.TagInfo
		err := tag.Update(c, db)
		if err != nil {
			cuserr.GinErrorHandle(err, c)
			return
		}
		c.JSON(200, gin.H{
			"message": "success",
		})
	})

	route.DELETE("/tag", func(c *gin.Context) {
		type Body struct {
			ID uint `form:"tag_id" binding:"required"`
		}
		var body Body
		err := c.ShouldBindQuery(&body)
		if err != nil {
			cuserr.GinErrorHandle(err, c)
		}
		db.Where("BlogID=?", body.ID).Delete(new(model.BlogTag))
		db.Where("ID=?", body.ID).Delete(new(model.TagInfo))
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})
}
