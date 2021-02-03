package model

type FileData struct {
	// gorm.Model
	ID          uint     `gorm:"primaryKey" json:"file_id"`
	Filename    string   `gorm:"not null;" json:"file_name"`
	FileHash    string   `gorm:"type:VARCHAR(512)" json:"-"`
	UserID      uint     `gorm:"not null;index" json:"user_id"`
	UserData    UserData `gorm:"not null;foreignKey:ID" json:"-"`
	BlogID      uint     `gorm:"not null;index" json:"blog_id"`
	ContextType string   `gorm:"not null" json:"Context-Type"`
	Deleted     uint     `gorm:"default:0" json:"-"` // 0: normal, 1: deleted
}
