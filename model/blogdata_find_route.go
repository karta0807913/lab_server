package model

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// this file generate by go generate, please don't edit it
// search options will put into struct
func (item *BlogData) Find(c *gin.Context, db *gorm.DB) ([]BlogData, error) {
	type Body struct {
		Title   *string `form:"title"`
		OwnerID *uint   `form:"user_id"`
	}
	var body Body
	var err error
	_ = c.ShouldBindQuery(&body)

	whereField := []string{
		"blog_data.deleted=?",
	}
	valueField := []interface{}{
		item.Deleted,
	}

	if body.Title != nil {
		whereField = append(whereField, "blog_data.title=?")
		valueField = append(valueField, body.Title)
		item.Title = *body.Title
	}

	if body.OwnerID != nil {
		whereField = append(whereField, "blog_data.owner_id=?")
		valueField = append(valueField, body.OwnerID)
		item.OwnerID = *body.OwnerID
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
	var result []BlogData
	if len(whereField) != 0 {
		db = db.Where(
			strings.Join(whereField, " and "),
			valueField[0], valueField[1:],
		)
	}
	err = db.Limit(limit).Offset(offset).Find(&result).Error
	return result, err
}
