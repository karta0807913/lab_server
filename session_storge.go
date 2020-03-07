package main

type Session interface {
	Get(string) interface{}
	Set(string, interface{}) error
	GetId() string
	SetId(string)
	Del(string)
	Clear()
}

type Storage interface {
	Get(string) (Session, error)
	Set(string, interface{}) error
	Create(interface{}) (Sessoin error)
	Del(string) error
}

type StorageError struct {
	error
	err_msg string
}

func (self *StorageError) Error() string {
	return self.err_msg
}

type StorageNotFoundError struct {
	StorageError
}

type SessionData struct {
	Session
	session map[string]interface{}
	updated bool
	id      string
}

func (self *SessionData) Clear() {
	self.session = make(map[string]interface{})
	self.updated = false
}

func (self *SessionData) SetId(key string) {
	self.id = key
}

func (self SessionData) GetId() string {
	return self.id
}

func (self SessionData) Get(key string) interface{} {
	return self.session[key]
}

func (self *SessionData) Set(key string, value interface{}) error {
	self.updated = true
	self.session[key] = value
	return nil
}

func (self *SessionData) Del(key string) {
	_, ok := self.session[key]
	self.updated = self.updated || ok
	delete(self.session, key)
}

func NewSessionData() *SessionData {
	return &SessionData{
		session: make(map[string]interface{}),
	}
}
