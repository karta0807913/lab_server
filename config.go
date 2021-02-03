package main

import (
	"os"
	"strconv"
)

type ConfigType struct {
	googleFileParent string
	sql              sqlConfig
	serverAddr       string
	publicKeyPath    string
	privateKeyPath   string
	googleAuthFile   string
	uploadPath       string
}

type mysqlConfig struct {
	sqlConfig
	account  string
	password string
	database string
	host     string
	port     int
}

type sqliteConfig struct {
	sqlConfig
	filepath string
}

type sqlConfig interface {
}

func number(varname string) int {
	i, err := strconv.Atoi(os.Getenv(varname))
	if err != nil {
		panic(err)
	}
	return i
}

var Config ConfigType = ConfigType{
	googleFileParent: "1376hSupEtrCgFDNmstrdK3oSDPOw5IGu",
	// sql: sqlConfig{
	// 	account:  os.Getenv("MYSQL_ACCOUNT"),
	// 	password: os.Getenv("MYSQL_PASS"),
	// 	database: os.Getenv("MYSQL_DATABS"),
	// 	host:     os.Getenv("MYSQL_HOST"),
	// 	port:     number("MYSQL_PORT"),
	// },
	// sql: sqlConfig{
	// 	account:  "test",
	// 	password: "123456",
	// 	database: "web_service",
	// 	host:     "172.18.0.2",
	// 	port:     3306,
	// },
	sql: sqliteConfig{
		filepath: "./sqlite.db",
	},
	publicKeyPath:  "./public.pem",
	privateKeyPath: "./private.pem",
	serverAddr:     ":1200",
	googleAuthFile: "./e539-lab-web-dd38239bcca2.json",
	uploadPath:     "./files",
}
