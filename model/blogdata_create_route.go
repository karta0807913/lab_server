package model

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// this file generate by go generate, please don't edit it
// data will put into struct
func (insert *BlogData) Create(c *gin.Context, db *gorm.DB) error {
	type Body struct {
		Title   string `json:"title" binding:"required"`
		Context string `json:"context" binding:"required"`

		FileList *[]FileData `json:"file_list"`
		TagList  *[]BlogTag  `json:"tag_list"`
	}
	var body Body
	err := c.ShouldBindJSON(&body)
	if err != nil {
		return err
	}

	selectField := []string{
		"title",
		"context",
		"owner_id",
		"created_at",
		"updated_at",
	}

	if body.FileList != nil {
		selectField = append(selectField, "file_list")
		insert.FileList = body.FileList
	}

	if body.TagList != nil {
		selectField = append(selectField, "tag_list")
		insert.TagList = body.TagList
	}

	insert.Title = body.Title
	insert.Context = body.Context

	return db.Select(
		selectField[0], selectField[1:],
	).Create(&insert).Error
}
