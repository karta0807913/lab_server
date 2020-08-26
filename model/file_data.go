package model

type FileData struct {
	// gorm.Model
	ID          uint     `gorm:"primaryKey" json:"file_id"`
	Filename    string   `gorm:"not null;" json:"file_name"`
	FileHash    string   `gorm:"type:VARCHAR(512)" json:"-"`
	UserId      uint     `gorm:"not null" json:"user_id"`
	ContextType string   `gorm:"not null" json:"Context-Type"`
	UserData    UserData `gorm:"not null;foreignKey:UserId" json:"-"`
	Deleted     uint     `gorm:"default:0" json:"-"` // 0: normal, 1: deleted
}
