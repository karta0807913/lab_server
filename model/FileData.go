package model

//go:generate generate_router -type "FileData" -ignore "UserData,ContextType" -method "Find"
type FileData struct {
	// gorm.Model
	ID          uint     `gorm:"primaryKey" json:"file_id"`
	FileName    string   `gorm:"not null;" json:"file_name"`
	FileHash    string   `gorm:"type:VARCHAR(512)" json:"-"`
	UserID      uint     `gorm:"not null;index" json:"user_id"`
	UserData    UserData `gorm:"not null;foreignKey:ID" json:"user_data"`
	BlogID      uint     `gorm:"not null;index" json:"blog_id"`
	ContextType string   `gorm:"not null" json:"Context-Type"`
	Deleted     uint     `gorm:"default:0;index" json:"deleted"` // 0: normal, 1: deleted
}
