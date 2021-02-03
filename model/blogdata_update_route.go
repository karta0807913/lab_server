package model

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// this file generate by go generate, please don't edit it
// data will put into struct
func (insert *BlogData) Update(c *gin.Context, db *gorm.DB) error {
	type Body struct {
		ID uint `json:"blog_id" binding:"required"`

		Title     *string     `json:"title"`
		OwnerID   *uint       `json:"user_id"`
		Owner     *UserData   `json:"owner"`
		Context   *string     `json:"context"`
		FileList  *[]FileData `json:"file_list"`
		TagList   *[]BlogTag  `json:"tag_list"`
		Deleted   *uint       `json:"deleted"`
		CreatedAt *time.Time  `json:"create_time"`
	}
	var body Body
	err := c.ShouldBindJSON(&body)
	if err != nil {
		return err
	}
	insert.ID = body.ID

	selectField := []string{
		"updated_at",
	}

	if body.Title != nil {
		selectField = append(selectField, "title")
		insert.Title = *body.Title
	}

	if body.OwnerID != nil {
		selectField = append(selectField, "owner_id")
		insert.OwnerID = *body.OwnerID
	}

	if body.Owner != nil {
		selectField = append(selectField, "owner")
		insert.Owner = body.Owner
	}

	if body.Context != nil {
		selectField = append(selectField, "context")
		insert.Context = *body.Context
	}

	if body.FileList != nil {
		selectField = append(selectField, "file_list")
		insert.FileList = body.FileList
	}

	if body.TagList != nil {
		selectField = append(selectField, "tag_list")
		insert.TagList = body.TagList
	}

	if body.Deleted != nil {
		selectField = append(selectField, "deleted")
		insert.Deleted = *body.Deleted
	}

	if body.CreatedAt != nil {
		selectField = append(selectField, "created_at")
		insert.CreatedAt = *body.CreatedAt
	}

	if len(selectField) == (0 + 1 + 1) {
		return errors.New("require at least one option")
	}

	return db.Select(
		selectField[0], selectField[1:],
	).Where("blog_data.id=?", body.ID).Updates(&insert).Error
}
