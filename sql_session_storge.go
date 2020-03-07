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
	db         *sql.DB
	get_sql    *sql.Stmt
	set_sql    *sql.Stmt
	insert_sql *sql.Stmt
	del_sql    *sql.Stmt
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
	insert_sql, err := db.Prepare(fmt.Sprintf("insert into `%s` (`data`) values (?)", config.TableName))
	if err != nil {
		return nil, err
	}
	del_sql, err := db.Prepare(fmt.Sprintf("delete from %s where id = ?", config.TableName))
	if err != nil {
		return nil, err
	}
	storage := &SQLStorage{
		db:         db,
		get_sql:    get_sql,
		set_sql:    set_sql,
		insert_sql: insert_sql,
		del_sql:    del_sql,
	}
	return storage, nil
}

func (self *SQLStorage) Get(session_id string) (Session, error) {
	result, err := self.get_sql.Query(session_id)
	if err != nil {
		return nil, err
	}
	if !result.Next() {
		return nil, &StorageNotFoundError{}
	}
	columns, err := result.Columns()
	if err != nil {
		return nil, err
	}
	body := make(map[string]interface{})
	err = json.Unmarshal([]byte(columns[0]), body)
	if err != nil {
		return nil, err
	}
	return &SessionData{
		session: body,
	}, nil
}

func (self *SQLStorage) Set(session_id string, body interface{}) error {
	var data []byte
	var err error
	switch session := body.(type) {
	case *SessionData:
	case SessionData:
		data, err = json.Marshal(session.session)
		break
	default:
		data, err = json.Marshal(session)
		break
	}
	if err != nil {
		return err
	}
	if session_id == "" {
		result, err := self.insert_sql.Exec(data)
		if err != nil {
			return err
		}
		id, err := result.LastInsertId()
		sid := fmt.Sprintf("%d", id)
		switch session := body.(type) {
		case Session:
			session.SetId(sid)
			break
		default:
			break
		}
	} else {
		_, err = self.set_sql.Exec(session_id, data, data)
		return err
	}
	return nil
}

func (self *SQLStorage) Del(session_id string) error {
	_, err := self.del_sql.Exec(session_id)
	return err
}
