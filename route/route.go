package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/karta0807913/go_server_utils/serverutil"
	"gorm.io/gorm"
)

type RouteConfig struct {
	Server     *gin.Engine
	UploadPath string
	DB         *gorm.DB
}

func checkLogin(c *gin.Context) {
	session := c.MustGet("session").(serverutil.Session)
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
		c.Header("Access-Control-Allow-Headers", "content-type,x-requested-with")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
	})

	apiRoute := config.Server.Group("/api/")
	APIRouteRegisterHandler(APIRouteConfig{
		db:    config.DB,
		route: apiRoute,
	})

	fileRoute := config.Server.Group("/file/")
	fileRoute.Use(checkLogin)
	FileRouteRegisterHandler(FileRouteConfig{
		route:      fileRoute,
		db:         config.DB,
		uploadPath: config.UploadPath,
	})

	memberRoute := config.Server.Group("/member/")
	memberRoute.Use(checkLogin)
	MemberRouteRegisterHandler(MemberRouteConfig{
		route: memberRoute,
		db:    config.DB,
	})

	blogRoute := config.Server.Group("/blog")
	BlogRouteRegisterHandler(APIRouteConfig{
		route: blogRoute,
		db:    config.DB,
	})

	websiteRoute := config.Server.Group("/web")
	WebsiteRouteRegisterHandler(WebsiteRouteConfig{
		route:    websiteRoute,
		prefix:   "/web",
		servPath: "./build",
	})

	adminRoute := config.Server.Group("/admin")
	AdminRouteRegisterHandler(APIRouteConfig{
		route: adminRoute,
		db:    config.DB,
	})
}
