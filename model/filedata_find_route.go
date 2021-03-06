package model

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// this file generate by go generate, please don't edit it
// search options will put into struct
func (item *FileData) Find(c *gin.Context, db *gorm.DB) ([]FileData, error) {
	type Body struct {
		ID      *uint `form:"file_id"`
		UserID  *uint `form:"user_id"`
		BlogID  *uint `form:"blog_id"`
		Deleted *uint `form:"deleted"`
	}
	var body Body
	var err error
	_ = c.ShouldBindQuery(&body)

	whereField := make([]string, 0)
	valueField := make([]interface{}, 0)

	if body.ID != nil {
		whereField = append(whereField, "file_data.id=?")
		valueField = append(valueField, body.ID)
		item.ID = *body.ID
	}

	if body.UserID != nil {
		whereField = append(whereField, "file_data.user_id=?")
		valueField = append(valueField, body.UserID)
		item.UserID = *body.UserID
	}

	if body.BlogID != nil {
		whereField = append(whereField, "file_data.blog_id=?")
		valueField = append(valueField, body.BlogID)
		item.BlogID = *body.BlogID
	}

	if body.Deleted != nil {
		whereField = append(whereField, "file_data.deleted=?")
		valueField = append(valueField, body.Deleted)
		item.Deleted = *body.Deleted
	}

	var limit int = 20
	slimit, ok := c.GetQuery("limit")
	if ok {
		limit, err = strconv.Atoi(slimit)
		if err != nil {
			limit = 20
		} else {
			if limit <= 0 || 20 < limit {
				limit = 20
			}
		}
	}
	soffset, ok := c.GetQuery("offset")
	var offset int
	if ok {
		offset, err = strconv.Atoi(soffset)
		if err != nil {
			offset = 0
		} else if offset < 0 {
			offset = 0
		}
	} else {
		offset = 0
	}
	var result []FileData
	if len(whereField) != 0 {
		db = db.Where(
			strings.Join(whereField, " and "),
			valueField[0], valueField[1:],
		)
	}
	err = db.Limit(limit).Offset(offset).Find(&result).Error
	return result, err
}
