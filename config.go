package main

type ConfigType struct {
	google_file_parent string
	sql                sqlConfig
	serverAddr         string
	public_key_path    string
	private_key_path   string
	google_auth_file   string
}

type sqlConfig struct {
	account  string
	password string
	database string
	host     string
	port     int
}

var Config ConfigType = ConfigType{
	google_file_parent: "1376hSupEtrCgFDNmstrdK3oSDPOw5IGu",
	sql: sqlConfig{
		account:  "e539",
		password: "e539lab",
		database: "web_service",
		host:     "127.0.0.1",
		port:     3306,
	},
	public_key_path:  "./public.pem",
	private_key_path: "./private.pem",
	serverAddr:       ":1200",
	google_auth_file: "./e539-lab-web-dd38239bcca2.json",
}
