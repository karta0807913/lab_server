package main

import (
	"os"
	"strconv"
)

type ConfigType struct {
	google_file_parent string
	sql                sqlConfig
	serverAddr         string
	public_key_path    string
	private_key_path   string
	google_auth_file   string
	upload_path        string
}

type sqlConfig struct {
	account  string
	password string
	database string
	host     string
	port     int
}

func number(varname string) int {
	i, err := strconv.Atoi(os.Getenv(varname))
	if err != nil {
		panic(err)
	}
	return i
}

var Config ConfigType = ConfigType{
	google_file_parent: "1376hSupEtrCgFDNmstrdK3oSDPOw5IGu",
	sql: sqlConfig{
		account:  os.Getenv("MYSQL_ACCOUNT"),
		password: os.Getenv("MYSQL_PASS"),
		database: os.Getenv("MYSQL_DATABS"),
		host:     os.Getenv("MYSQL_HOST"),
		port:     number("MYSQL_PORT"),
	},
	public_key_path:  "./public.pem",
	private_key_path: "./private.pem",
	serverAddr:       ":1200",
	google_auth_file: "./e539-lab-web-dd38239bcca2.json",
	upload_path:      "./files",
}
