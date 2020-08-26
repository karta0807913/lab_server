package server

import (
	"encoding/json"
	"log"
	"strconv"

	"gorm.io/gorm"
)

type GormStorage struct {
	Storage
	db *gorm.DB
}

type SessionModel struct {
	gorm.Model
	ID   uint   `gorm:"primaryKey"`
	Data []byte `gorm:"type:json"`
}

func NewGormStorage(db *gorm.DB) (*GormStorage, error) {
	err := db.AutoMigrate(&SessionModel{})
	if err != nil {
		return nil, err
	}
	return &GormStorage{
		db: db,
	}, nil
}

func (self GormStorage) Get(session_id string) (Session, error) {
	var session SessionModel
	tx := self.db.First(&session, "id = ?", session_id)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if session.Data == nil {
		return nil, &StorageNotFoundError{}
	}
	body := make(map[string]interface{})
	err := json.Unmarshal(session.Data, &body)
	if err != nil {
		return nil, err
	}
	id := strconv.Itoa(int(session.ID))
	return &MapSession{
		session: body,
		id:      id,
	}, nil
}

func (self *GormStorage) Set(session Session) error {
	data, err := json.Marshal(session.All())
	if err != nil {
		return err
	}
	model := &SessionModel{
		Data: data,
	}
	if session.GetId() != "" {
		sid, err := strconv.Atoi(session.GetId())
		if err != nil {
			return err
		}
		model.ID = uint(sid)
	}
	tx := self.db.Create(model)
	if tx.Error != nil {
		tx = self.db.Model(model).Updates(model)
		if tx.Error != nil {
			log.Println(tx.Error)
			return tx.Error
		}
	}
	session.SetId(strconv.Itoa(int(model.ID)))
	return nil
}

func (self *GormStorage) Del(session_id string) error {
	tx := self.db.Delete(&SessionModel{}, "id = ?", session_id)
	return tx.Error
}
