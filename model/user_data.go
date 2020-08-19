package model

type UserData struct {
	// gorm.Model
	ID       uint `gorm:"primaryKey"`
	Nickname string
	Account  string
	Password string
	Status   uint `gorm:"default:0"` // 0: normal user
}
