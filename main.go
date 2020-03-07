package main

import (
	"database/sql"
	"fmt"
	"log"
)

func clearDrive(drive *GoogleDrive) {
	files, err := drive.ListFiles()
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files.Files {
		if file.Id == "0Bzaq6TKWYiNjc3RhcnRlcl9maWxl" || file.Id == "1376hSupEtrCgFDNmstrdK3oSDPOw5IGu" {
			continue
		}
		drive.DeleteFile(file.Id)
	}
}

func main() {
	db, err := sql.Open(
		"mysql",
		fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s",
			Config.sql.account,
			Config.sql.password,
			Config.sql.host,
			Config.sql.port,
			Config.sql.database,
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	storage, err := NewSQLStorage(db, SQLStorageConfig{
		TableName: "web_session",
	})
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	token, err := google_token(
		Config.google_auth_file,
		"https://www.googleapis.com/auth/drive",
	)
	if err != nil {
		log.Fatal(err)
	}
	drive := &GoogleDrive{
		token: token,
	}

	server, err :=
		NewSessionHttpServer(ServerSettings{
			PublicKeyPath:  Config.public_key_path,
			PrivateKeyPath: Config.private_key_path,
			ServerAddress:  Config.serverAddr,
			Db:             db,
			Storage:        storage,
			Drive:          drive,
		})
	if err != nil {
		log.Fatal(err)
	}

	route(server)

	log.Printf("server listening on %s", Config.serverAddr)
	log.Fatal(server.ListenAndServe())
}
