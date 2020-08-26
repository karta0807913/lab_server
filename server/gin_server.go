package server

import (
	"github.com/gin-gonic/gin"
	"github.com/karta0807913/lab_server/utils"
)

func NewGinServer(config ServerSettings) (*gin.Engine, error) {
	jwt, err := utils.NewJwtHelper(config.PublicKeyPath, config.PrivateKeyPath)
	if err != nil {
		return nil, err
	}
	server := gin.New()
	server.Use(NewGinSessionFactory(jwt, config.Storage).SessionMiddleware(config.SessionName))
	return server, nil
}
