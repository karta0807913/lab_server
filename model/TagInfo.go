package model

//go:generate generate_router -type "TagInfo" -method "Create"
//go:generate generate_router -type "TagInfo" -method "First"
//go:generate generate_router -type "TagInfo" -method "Update" -require "Name"
type TagInfo struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"not null" json:"name"`
}
