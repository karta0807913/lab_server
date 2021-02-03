package model

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// this file generate by go generate, please don't edit it
// data will put into struct
func (insert *BlogTag) Create(c *gin.Context, db *gorm.DB) error {
	type Body struct {
		BlogID uint `json:"blog_id" binding:"required"`
		TagID  uint `json:"tag_id" binding:"required"`
	}
	var body Body
	err := c.ShouldBindJSON(&body)
	if err != nil {
		return err
	}

	selectField := []string{
		"blog_id",
		"tag_id",
	}

	insert.BlogID = body.BlogID
	insert.TagID = body.TagID

	return db.Select(
		selectField[0], selectField[1:],
	).Create(&insert).Error
}
