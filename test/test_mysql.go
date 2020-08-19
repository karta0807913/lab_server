package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/karta0807913/lab_server/model"
)

type User struct {
	A *string `json:"a"`
	B *string `json:"b"`
}

func TestMySql(t *testing.T) {
	var user User
	err := json.Unmarshal([]byte(`{"a": "1"}`), &user)
	fmt.Println(err)
	fmt.Printf("%+v\n", user)
	model.CreateDB("test", "123456", "172.18.0.2", 3306, "test")
}
