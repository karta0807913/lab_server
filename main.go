package main

import (
	"fmt"
	"log"

	"github.com/karta0807913/go_server_utils/serverutil"
	"github.com/karta0807913/lab_server/model"
	"github.com/karta0807913/lab_server/route"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var db *gorm.DB
	var err error
	switch sql := Config.sql.(type) {
	case mysqlConfig:
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			sql.account,
			sql.password,
			sql.host,
			sql.port,
			sql.database,
		)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal(err)
		}
	case sqliteConfig:
		db, err = model.CreateSqliteDB(sql.filepath)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = model.InitDB(db)
	if err != nil {
		log.Fatalln(err)
	}
	storage, err := serverutil.NewGormStorage(db)
	if err != nil {
		log.Fatal(err)
	}

	server, err :=
		serverutil.NewGinServer(serverutil.ServerSettings{
			PrivateKeyPath: Config.privateKeyPath,
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
		UploadPath: Config.uploadPath,
	})

	log.Printf("server listening on %s", Config.serverAddr)
	log.Fatal(server.Run(Config.serverAddr))
}
