package main

//
// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// )
//
// type AA struct {
// 	A string `json:"a"`
// 	C string `json:"c"`
// }
//
// func main() {
// 	decoder := json.NewDecoder(bytes.NewReader([]byte(`{ "a": "b", "c": "d" }`)))
// 	decoder.DisallowUnknownFields()
// 	bb := new(AA)
// 	if err := decoder.Decode(bb); err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println(bb)
// }
//
