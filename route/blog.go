package route

import (
	"fmt"
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
		err := data.First(c, db.Preload("FileList").Preload("TagList.TagInfo").Preload("Owner").Select("ID", "Title", "OwnerID", "Context", "CreatedAt", "UpdatedAt"))
		fmt.Println(data.TagList)
		if err != nil {
			cuserr.GinErrorHandle(err, c)
			return
		}
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
		db.Joins("BlogData").Preload("BlogData").Preload("BlogData.TagList.TagInfo").Preload("BlogData.Owner").Where("BlogData.Deleted=0 and tag_id in ?", tagList).Order("blog_id desc").Group("blog_id").Find(&result)
		c.JSON(200, result)
	})

	route.POST("/blog_tag", func(c *gin.Context) {
		var blogTag model.BlogTag
		err := blogTag.Create(c, db)
		if err != nil {
			cuserr.GinErrorHandle(err, c)
			return
		}
		c.JSON(200, blogTag)
	})

	// TODO: permission check
	route.DELETE("/blog_tag", func(c *gin.Context) {
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
		err = db.Where("blog_id=? and tag_id=?", body.BlogID, body.TagID).Delete(new(model.BlogTag)).Error
		if err != nil {
			cuserr.GinErrorHandle(err, c)
			return
		}
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})

	route.POST("/blog", func(c *gin.Context) {
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

	route.PUT("/blog", func(c *gin.Context) {
		var blog model.BlogData
		blog.UpdatedAt = time.Now()
		err := blog.Update(c, db)
		if err != nil {
			cuserr.GinErrorHandle(err, c)
			return
		}
		c.JSON(200, gin.H{
			"message": "success",
		})
	})

	route.DELETE("/blog", func(c *gin.Context) {
		type Body struct {
			ID uint `json:"blog_id" binding:"required"`
		}
		var body Body
		err := c.ShouldBindJSON(&body)
		if err != nil {
			cuserr.GinErrorHandle(err, c)
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

	route.POST("/tag", func(c *gin.Context) {
		var tagInfo model.TagInfo
		err := tagInfo.Create(c, db)
		if err != nil {
			cuserr.GinErrorHandle(err, c)
			return
		}
		c.JSON(200, gin.H{"message": "ok"})
	})

	route.PUT("/tag", func(c *gin.Context) {
		var tagInfo model.TagInfo
		err := tagInfo.Update(c, db)
		if err != nil {
			cuserr.GinErrorHandle(err, c)
			return
		}
		c.JSON(200, gin.H{"message": "ok"})
	})
}
