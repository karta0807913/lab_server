package model

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (insert *UserData) Update(c *gin.Context, db *gorm.DB) error {
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
		return err
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
		insert.Password = *body.Password
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
		return errors.New("require at least one option")
	}

	return db.Select(
		selectField[0], selectField[1:],
	).Where("user_data.id=?", body.ID).Updates(&insert).Error
}
