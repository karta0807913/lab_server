package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
)

func saltPassword(str string) string {
	encoder := sha256.New()
	password := base64.StdEncoding.EncodeToString(
		encoder.Sum(
			[]byte(str),
		),
	)
	return password
}

func main() {
	fmt.Println(saltPassword(os.Args[1]))
}
