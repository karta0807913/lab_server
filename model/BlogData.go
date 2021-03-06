package model

import "time"

//go:generate generate_router -type "BlogData" -method "Create" -ignore "Deleted,Owner" -options "FileList,TagList" -useDefault "OwnerID,CreatedAt,UpdatedAt" -minItem 0
//go:generate generate_router -type "BlogData" -method "First" -ignore "Title,OwnerID" -require "ID"
//go:generate generate_router -type "BlogData" -method "Update" -useDefault "UpdatedAt"
//go:generate generate_router -type "BlogData" -ignore "Owner,FileList,TagList,ID" -useDefault "Deleted" -method "Find"
type BlogData struct {
	ID        uint        `gorm:"primaryKey" json:"blog_id"`
	Title     string      `gorm:"index;not null;type:VARCHAR(128)" json:"title"`
	OwnerID   uint        `gorm:"index; not null" json:"user_id"`
	Owner     *UserData   `gorm:"not null;foreignKey:ID;references:OwnerID" json:"owner"`
	Context   string      `gorm:"not null;type:LONGTEXT" json:"context"`
	FileList  *[]FileData `gorm:"foreignKey:BlogID" json:"file_list"`
	TagList   *[]BlogTag  `gorm:"foreignKey:BlogID" json:"tag_list"`
	Deleted   uint        `gorm:"not null;default:0;index" json:"deleted"`
	CreatedAt time.Time   `gorm:"not null" json:"create_time"`
	UpdatedAt time.Time   `gorm:"not null" json:"update_time"`
}
