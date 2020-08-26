package route

import (
	"os"
	"path"

	"github.com/gin-gonic/gin"
)

type WebsiteRouteConfig struct {
	route    *gin.RouterGroup
	prefix   string
	servPath string
}

func WebsiteRouteRegistHandler(config WebsiteRouteConfig) {
	route := config.route

	route.GET("/*name", func(c *gin.Context) {
		name := c.Param("name")[1:]
		filePath := path.Join(config.servPath, name)
		_, err := os.Stat(filePath)
		if err != nil {
			c.File(path.Join(config.servPath, "index.html"))
		} else {
			c.File(path.Join(config.servPath, name))
		}
	})
}
