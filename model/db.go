package model

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(db *gorm.DB) error {
	err := db.AutoMigrate(&UserData{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&FileData{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&TagInfo{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&BlogData{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&BlogTag{})
	if err != nil {
		return err
	}
	return nil
}
func CreateSqliteDB(filename string) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(filename), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   logger.Default.LogMode(logger.Info),
	})
}
