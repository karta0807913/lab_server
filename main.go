package main

import (
	"log"

	"github.com/karta0807913/lab_server/model"
	"github.com/karta0807913/lab_server/route"
	"github.com/karta0807913/lab_server/server"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
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

	storage, err := server.NewGormStorage(db)
	if err != nil {
		log.Fatal(err)
	}

	server, err :=
		server.NewGinServer(server.ServerSettings{
			PublicKeyPath:  Config.public_key_path,
			PrivateKeyPath: Config.private_key_path,
			ServerAddress:  Config.serverAddr,
			Db:             db,
			Storage:        storage,
			SessionName:    "session",
		})
	if err != nil {
		log.Fatal(err)
	}

	route.Route(route.RouteConfig{
		DB:         db,
		Server:     server,
		UploadPath: Config.upload_path,
	})

	log.Printf("server listening on %s", Config.serverAddr)
	log.Fatal(server.Run(Config.serverAddr))
}
