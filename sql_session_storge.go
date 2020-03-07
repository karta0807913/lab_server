package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type SQLStorageConfig struct {
	TableName string
}

type SQLStorage struct {
	Storage
	db      *sql.DB
	get_sql *sql.Stmt
	set_sql *sql.Stmt
	del_sql *sql.Stmt
}

func NewSQLStorage(db *sql.DB, config SQLStorageConfig) (*SQLStorage, error) {
	get_sql, err := db.Prepare(fmt.Sprintf("select data from `%s` where id=? limit 1", config.TableName))
	if err != nil {
		return nil, err
	}
	set_sql, err := db.Prepare(fmt.Sprintf("insert into `%s` (`id`, `data`) values (?, ?) ON DUPLICATE KEY UPDATE data=?", config.TableName))
	if err != nil {
		return nil, err
	}
	del_sql, err := db.Prepare(fmt.Sprintf("delete from %s where id = ?", config.TableName))
	if err != nil {
		return nil, err
	}
	storage := &SQLStorage{
		db:      db,
		get_sql: get_sql,
		set_sql: set_sql,
		del_sql: del_sql,
	}
	return storage, nil
}

func (self *SQLStorage) Get(session_id string, body interface{}) error {
	result, err := self.get_sql.Query(session_id)
	if err != nil {
		return err
	}
	if !result.Next() {
		return &StorageNotFoundError{}
	}
	columns, err := result.Columns()
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(columns[0]), body)
	return err
}

func (self *SQLStorage) Set(session_id string, body interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}
	_, err = self.set_sql.Exec(session_id, data, data)
	return err
}

func (self *SQLStorage) Del(session_id string) error {
	_, err := self.del_sql.Exec(session_id)
	return err
}
