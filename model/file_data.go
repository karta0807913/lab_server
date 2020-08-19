package model

type FileData struct {
	// gorm.Model
	ID       uint `gorm:"primaryKey`
	Filename string
	FileHash string `gorm:"type:VARCHAR(512)"`
	Deleted  uint   `gorm:"default:0"` // 0: normal, 1: deleted
}
