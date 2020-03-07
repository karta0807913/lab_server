package main

type Storage interface {
	Get(string, interface{}) error
	Set(string, interface{}) error
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
