package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/karta0807913/lab_server/server"
	"gorm.io/gorm"
)

type RouteConfig struct {
	Server     *gin.Engine
	UploadPath string
	DB         *gorm.DB
}

func checkLogin(c *gin.Context) {
	session := c.MustGet("session").(server.Session)
	id := session.Get("mem_id")
	if id == nil {
		c.String(403, "Permission Denied")
		c.Abort()
	}
}

func Route(config RouteConfig) {
	// serv.GET("/", func(c *gin.Context) {
	// 	session := c.MustGet("session").(server.Session)
	// 	fmt.Println(session.Get("A"))
	// 	if session.Get("A") != nil && session.Get("A").(string) == "B" {
	// 		session.Set("A", "C")
	// 	} else {
	// 		session.Set("A", "B")
	// 	}
	// 	c.Writer.Write([]byte("Hello"))
	// })
	config.Server.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://127.0.0.1:3000")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "content-type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
	})

	api_route := config.Server.Group("/api/")
	ApiRouteRegistHandler(ApiRouteConfig{
		db:    config.DB,
		route: api_route,
	})

	file_route := config.Server.Group("/file/")
	file_route.Use(checkLogin)
	FileRouteRegistHandler(FileRouteConfig{
		route:      file_route,
		db:         config.DB,
		uploadPath: config.UploadPath,
	})

	member_route := config.Server.Group("/member/")
	member_route.Use(checkLogin)
	MemberRouteRegistHandler(MemberRouteConfig{
		route: member_route,
		db:    config.DB,
	})

	website_route := config.Server.Group("/web")
	WebsiteRouteRegistHandler(WebsiteRouteConfig{
		route:    website_route,
		prefix:   "/web",
		servPath: "./build",
	})
}
