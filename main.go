package main

import (
	"log"

	"github.com/karta0807913/lab_server/model"
	"github.com/karta0807913/lab_server/route"
	"github.com/karta0807913/lab_server/server"
)

func main() {
	db, err := model.CreateDB(
		Config.sql.account,
		Config.sql.password,
		Config.sql.host,
		Config.sql.port,
		Config.sql.database,
	)
	if err != nil {
		log.Fatal(err)
	}

	sqlDb, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	storage, err := server.NewSQLStorage(sqlDb, server.SQLStorageConfig{
		TableName: "web_session",
	})
	if err != nil {
		log.Fatal(err)
	}

	server, err :=
		server.NewSessionHttpServer(server.ServerSettings{
			PublicKeyPath:  Config.public_key_path,
			PrivateKeyPath: Config.private_key_path,
			ServerAddress:  Config.serverAddr,
			Db:             db,
			Storage:        storage,
		})
	if err != nil {
		log.Fatal(err)
	}

	route.Route(server, Config.upload_path)

	log.Printf("server listening on %s", Config.serverAddr)
	log.Fatal(server.ListenAndServe())
}
