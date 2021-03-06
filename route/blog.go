package route

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/karta0807913/go_server_utils/serverutil"
	cuserr "github.com/karta0807913/lab_server/error"
	"github.com/karta0807913/lab_server/model"
)

func BlogRouteRegisterHandler(config APIRouteConfig) {
	db := config.db
	route := config.route

	AdminCheck := func(c *gin.Context) bool {
		session := c.MustGet("session").(serverutil.Session)
		return session.Get("is_admin").(bool)
	}

	BlogOwnerCheck := func(BlogID uint, c *gin.Context) bool {
		var data model.BlogData
		err := db.Select("OwnerID").Where("id=?", BlogID).First(&data).Error
		if err != nil {
			return false
		}
		session := c.MustGet("session").(serverutil.Session)
		memID := uint(session.Get("mem_id").(float64))
		return data.OwnerID == memID
	}

	BlogTagOwnerCheck := func(BlogID uint, TagID uint, c *gin.Context) bool {
		var data model.BlogTag
		err := db.Preload("BlogData").Where("blog_id=? and tag_id=?", BlogID, TagID).First(&data).Error
		if err != nil {
			return false
		}
		session := c.MustGet("session").(serverutil.Session)
		memID := uint(session.Get("mem_id").(float64))
		return data.BlogData.OwnerID == memID
	}

	// AdminCheck = func(c *gin.Context) bool {
	// 	session := c.MustGet("session").(serverutil.Session)
	// 	var user model.UserData
	// 	mem_id := uint(session.Get("mem_id").(float64))

	// 	err := db.Select("IsAdmin").Where("mem_id=?", mem_id).First(&user).Error
	// 	if err != nil {
	// 		log.Println("get userdata %s failed", mem_id)
	// 		return false
	// 	}
	// 	return user.IsAdmin
	// }

	route.GET("/get", func(c *gin.Context) {
		var data model.BlogData
		err := data.First(c, db.Preload("FileList", "deleted=0").Preload("TagList.TagInfo").Preload("Owner").Select("ID", "Title", "OwnerID", "Context", "CreatedAt", "UpdatedAt"))
		if err != nil {
			cuserr.GinErrorHandle(err, c)
			return
		}
		data.Owner.IsAdmin = false
		c.JSON(200, data)
	})

	route.GET("/list", func(c *gin.Context) {
		var data model.BlogData
		data.Deleted = 0
		res, err := data.Find(
			c,
			db.Select("blog_data.id", "title", "OwnerID", "CreatedAt", "UpdatedAt").Preload("TagList.TagInfo").Order("blog_data.id desc"),
		)
		if err != nil {
			cuserr.GinErrorHandle(err, c)
			return
		}
		c.JSON(200, res)
	})

	route.GET("/blog_tag", func(c *gin.Context) {
		tagIDString, ok := c.GetQuery("tag_id")
		if !ok {
			c.JSON(403, gin.H{
				"message": "key tag_id not found",
			})
			return
		}
		tagIDArray := strings.Split(tagIDString, ",")
		tagList := make([]uint, 0, len(tagIDArray))
		for _, tagString := range tagIDArray {
			tagID, err := strconv.Atoi(tagString)
			if err != nil || tagID <= 0 {
				continue
			}
			tagList = append(tagList, uint(tagID))
		}

		var result []model.BlogTag
		db.Select(
			"blog_tags.id", "blog_id", "tag_id",
		).Joins(
			"inner join blog_data on blog_data.id = blog_id and deleted=0",
		).Preload(
			"BlogData",
		).Preload(
			"BlogData.TagList.TagInfo",
		).Preload(
			"BlogData.Owner",
		).Where(
			"tag_id in ?", tagList,
		).Order(
			"blog_id desc",
		).Group(
			"blog_tags.id",
		).Group(
			"blog_id",
		).Group(
			"tag_id",
		).Find(&result)
		c.JSON(200, result)
	})

	route.POST("/blog_tag", checkLogin, func(c *gin.Context) {
		type Body struct {
			BlogID uint `json:"blog_id" binding:"required"`
			TagID  uint `json:"tag_id" binding:"required"`
		}
		var insert model.BlogTag
		var body Body
		err := c.ShouldBindJSON(&body)
		if err != nil {
			cuserr.GinErrorHandle(err, c)
			return
		}

		if !AdminCheck(c) && !BlogOwnerCheck(body.BlogID, c) {
			c.JSON(403, gin.H{
				"message": "permission denied",
			})
			return
		}

		selectField := []string{
			"blog_id",
			"tag_id",
		}

		insert.BlogID = body.BlogID
		insert.TagID = body.TagID

		err = db.Select(
			selectField[0], selectField[1:],
		).Create(&insert).Error
		if err != nil {
			cuserr.GinErrorHandle(err, c)
			return
		}
		c.JSON(200, insert)
	})

	route.DELETE("/blog_tag", checkLogin, func(c *gin.Context) {
		type Body struct {
			TagID  uint `json:"tag_id"`
			BlogID uint `json:"blog_id"`
		}
		var body Body
		err := c.ShouldBindJSON(&body)
		if err != nil {
			cuserr.GinErrorHandle(err, c)
			return
		}
		if !AdminCheck(c) && !BlogTagOwnerCheck(body.BlogID, body.TagID, c) {
			c.JSON(403, gin.H{
				"message": "permission denied",
			})
			return
		}
		err = db.Where("blog_id=? and tag_id=?", body.BlogID, body.TagID).Delete(new(model.BlogTag)).Error
		if err != nil {
			cuserr.GinErrorHandle(err, c)
			return
		}
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})

	route.POST("/blog", checkLogin, func(c *gin.Context) {
		var blog model.BlogData
		session := c.MustGet("session").(serverutil.Session)
		blog.OwnerID = uint(session.Get("mem_id").(float64))
		blog.CreatedAt = time.Now()
		blog.UpdatedAt = time.Now()
		err := blog.Create(c, db)
		if err != nil {
			cuserr.GinErrorHandle(err, c)
			return
		}

		c.JSON(200, gin.H{
			"blog_id": blog.ID,
			"message": "success",
		})
	})

	route.PUT("/blog", checkLogin, func(c *gin.Context) {
		var blog model.BlogData
		blog.UpdatedAt = time.Now()
		type Body struct {
			ID uint `json:"blog_id" binding:"required"`

			Title     *string           `json:"title"`
			OwnerID   *uint             `json:"user_id"`
			Owner     *model.UserData   `json:"owner"`
			Context   *string           `json:"context"`
			FileList  *[]model.FileData `json:"file_list"`
			TagList   *[]model.BlogTag  `json:"tag_list"`
			Deleted   *uint             `json:"deleted"`
			CreatedAt *time.Time        `json:"create_time"`
		}
		var body Body
		err := c.ShouldBindJSON(&body)
		if err != nil {
			cuserr.GinErrorHandle(err, c)
			return
		}

		if !AdminCheck(c) && !BlogOwnerCheck(body.ID, c) {
			c.JSON(403, gin.H{
				"message": "permission denied",
			})
			return
		}

		blog.ID = body.ID

		selectField := []string{
			"updated_at",
		}

		if body.Title != nil {
			selectField = append(selectField, "title")
			blog.Title = *body.Title
		}

		if body.OwnerID != nil {
			selectField = append(selectField, "owner_id")
			blog.OwnerID = *body.OwnerID
		}

		if body.Owner != nil {
			selectField = append(selectField, "owner")
			blog.Owner = body.Owner
		}

		if body.Context != nil {
			selectField = append(selectField, "context")
			blog.Context = *body.Context
		}

		if body.FileList != nil {
			selectField = append(selectField, "file_list")
			blog.FileList = body.FileList
		}

		if body.TagList != nil {
			selectField = append(selectField, "tag_list")
			blog.TagList = body.TagList
		}

		if body.Deleted != nil {
			selectField = append(selectField, "deleted")
			blog.Deleted = *body.Deleted
		}

		if body.CreatedAt != nil {
			selectField = append(selectField, "created_at")
			blog.CreatedAt = *body.CreatedAt
		}

		if len(selectField) == (0 + 1 + 1) {
			c.JSON(403, gin.H{
				"message": "require at least one option",
			})
			return
		}

		err = db.Select(
			selectField[0], selectField[1:],
		).Where("blog_data.id=?", body.ID).Updates(&blog).Error
		if err != nil {
			cuserr.GinErrorHandle(err, c)
			return
		}
		c.JSON(200, gin.H{
			"message": "success",
		})
	})

	route.DELETE("/blog", checkLogin, func(c *gin.Context) {
		type Body struct {
			ID uint `json:"blog_id" binding:"required"`
		}
		var body Body
		err := c.ShouldBindJSON(&body)
		if err != nil {
			cuserr.GinErrorHandle(err, c)
			return
		}
		if !AdminCheck(c) && !BlogOwnerCheck(body.ID, c) {
			c.JSON(403, gin.H{
				"message": "permission denied",
			})
			return
		}
		err = db.Model(new(model.BlogData)).Where("id=?", body.ID).Update("Deleted", 1).Error
		if err != nil {
			cuserr.GinErrorHandle(err, c)
			return
		}
		c.JSON(200, gin.H{
			"message": "success",
		})
	})

	route.GET("/tag", func(c *gin.Context) {
		var tagInfo model.TagInfo
		result, _ := tagInfo.Find(c, db)
		c.JSON(200, result)
	})
}
