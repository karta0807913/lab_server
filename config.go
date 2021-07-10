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
	calendarID       string
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

func getSQLConfig() sqlConfig {
	var sql sqlConfig = sqliteConfig{
		filepath: "./sqlite.db",
	}

	if os.Getenv("FORMAL") == "true" {
		sql = mysqlConfig{
			account:  os.Getenv("MYSQL_ACCOUNT"),
			password: os.Getenv("MYSQL_PASS"),
			database: os.Getenv("MYSQL_DATABS"),
			host:     os.Getenv("MYSQL_HOST"),
			port:     number("MYSQL_PORT"),
		}
	}
	return sql
}

var WebsiteConfig ConfigType = ConfigType{
	googleFileParent: "1376hSupEtrCgFDNmstrdK3oSDPOw5IGu",
	sql:              getSQLConfig(),
	publicKeyPath:    "./public.pem",
	privateKeyPath:   "./private.pem",
	serverAddr:       ":1200",
	googleAuthFile:   "./e539-lab-web-9227fbd1854a.json",
	uploadPath:       "./files",
	calendarID:       "pp4f60hjm0llrf1pslqeoavho8@group.calendar.google.com",
}
