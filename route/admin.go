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
		err := c.ShouldBindJSON(&body)
		if err != nil {
			cuserr.GinErrorHandle(err, c)
			return
		}
		db.Where("TagID=?", body.ID).Delete(new(model.BlogTag))
		db.Where("ID=?", body.ID).Delete(new(model.TagInfo))
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})

	route.GET("/user", func(c *gin.Context) {
		var user model.UserData
		result, _ := user.Find(c, db.Select("ID, Nickname, Account, is_admin, Status"))
		c.JSON(200, result)
	})

	route.POST("/user", func(c *gin.Context) {
		var insert model.UserData
		type Body struct {
			Nickname string `json:"nickname" binding:"required"`
			Account  string `json:"account" binding:"required"`
			Password string `json:"password" binding:"required"`
			IsAdmin  bool   `json:"is_admin"`
			Status   uint   `json:"status"`
		}
		var body Body
		err := c.ShouldBindJSON(&body)
		if err != nil {
			cuserr.GinErrorHandle(err, c)
			return
		}

		selectField := []string{
			"nickname",
			"account",
			"password",
			"is_admin",
			"status",
		}

		insert.Nickname = body.Nickname
		insert.Account = body.Account
		insert.Password = saltPassword(body.Password)
		insert.IsAdmin = body.IsAdmin
		insert.Status = body.Status

		err = db.Select(
			selectField[0], selectField[1:],
		).Create(&insert).Error
		if err != nil {
			cuserr.GinErrorHandle(err, c)
			return
		}
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})

	route.PUT("/user", func(c *gin.Context) {
		var insert model.UserData
		type Body struct {
			ID uint `json:"user_id" binding:"required"`

			Nickname *string `json:"nickname"`
			Account  *string `json:"account"`
			Password *string `json:"password"`
			IsAdmin  *bool   `json:"is_admin"`
			Status   *uint   `json:"status"`
		}
		var body Body
		err := c.ShouldBindJSON(&body)
		if err != nil {
			cuserr.GinErrorHandle(err, c)
			return
		}
		insert.ID = body.ID

		selectField := make([]string, 0)

		if body.Nickname != nil {
			selectField = append(selectField, "nickname")
			insert.Nickname = *body.Nickname
		}

		if body.Account != nil {
			selectField = append(selectField, "account")
			insert.Account = *body.Account
		}

		if body.Password != nil {
			selectField = append(selectField, "password")
			insert.Password = saltPassword(*body.Password)
		}

		if body.IsAdmin != nil {
			selectField = append(selectField, "is_admin")
			insert.IsAdmin = *body.IsAdmin
		}

		if body.Status != nil {
			selectField = append(selectField, "status")
			insert.Status = *body.Status
		}

		if len(selectField) == (0 + 0 + 1) {
			c.JSON(403, gin.H{
				"message": "require at least one option",
			})
			return
		}

		err = db.Select(
			selectField[0], selectField[1:],
		).Where("user_data.id=?", body.ID).Updates(&insert).Error
		if err != nil {
			cuserr.GinErrorHandle(err, c)
			return
		}
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})
}
