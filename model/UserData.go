package model

//go:generate generate_router -type "UserData" -method "Find" -ignore "ID,Account"
type UserData struct {
	// gorm.Model
	ID       uint   `gorm:"primaryKey" json:"user_id"`
	Nickname string `gorm:"not null" json:"nickname"`
	Account  string `gorm:"uniqueIndex;not null;type:VARCHAR(15)" json:"-"`
	Password string `gorm:"not null" json:"-"`
	IsAdmin  bool   `gorm:"default:false" json:"is_admin"`
	Status   uint   `gorm:"default:0;not null" json:"-"`
	// 0: not active, 1: active, 2: deleted, not activated, 3 deleted and activated
	//   MSB        LSB
	// deleted   activated
	//    0          0
}
