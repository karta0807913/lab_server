package model

import (
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// this file generate by go generate, please don't edit it
// search options will put into struct
func (item *BlogData) First(c *gin.Context, db *gorm.DB) error {
	type Body struct {
		ID uint `form:"blog_id" binding:"required"`
	}

	var body Body
	err := c.ShouldBindQuery(&body)
	if err != nil {
		return err
	}

	whereField := []string{
		"blog_data.id=?",
	}
	valueField := []interface{}{
		body.ID,
	}

	item.ID = body.ID

	err = db.Where(
		strings.Join(whereField, " and "),
		valueField[0], valueField[1:],
	).First(item).Error
	return err
}
