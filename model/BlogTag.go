package model

//go:generate generate_router -type "BlogTag" -method "Create" -ignore "BlogData,TagInfo"
type BlogTag struct {
	ID       uint      `gorm:"primaryKey" json:"blog_tag_id"`
	BlogID   uint      `gorm:"index;not null" json:"blog_id"`
	BlogData *BlogData `gorm:"foreignKey:ID;references:BlogID" json:"blog_data"`
	TagID    uint      `gorm:"index;not null" json:"tag_id"`
	TagInfo  *TagInfo  `gorm:"foreignKey:ID;references:TagID" json:"tag_info"`
}
