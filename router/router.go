package router

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"go_restful_api/handler/sd"
	"go_restful_api/handler/user"
	"go_restful_api/router/middleware"
	_ "go_restful_api/docs"
	"net/http"
	"github.com/swaggo/gin-swagger" // gin-swagger middleware
	"github.com/swaggo/files"
)

func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	//middleware
	g.Use(gin.Recovery())
	g.Use(middleware.NoCache) //todo 影响性能，可以去掉
	g.Use(middleware.Options) //todo 影响性能，可以去掉
	g.Use(middleware.Secure)
	g.Use(mw...)

	//404 handler
	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The incorrect API route")
	})

	// swagger api docs
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// pprof router
	pprof.Register(g)

	// api for authentication functionalities
	g.POST("/login", user.Login)


	u := g.Group("/v1/user")
	u.Use(middleware.AuthMiddleware())
	{
		u.POST("", user.Create)
		u.DELETE("/:id", user.Delete)
		u.PUT("/:id", user.Update)
		u.GET("", user.List)
		u.GET("/:username", user.Get)
	}

	//health check handler
	svcd := g.Group("/sd")
	{
		svcd.GET("/health", sd.HealthCheck)
		svcd.GET("/disk", sd.DiskCheck)
		svcd.GET("/cpu", sd.CUPCheck)
		svcd.GET("/ram", sd.RAMCheck)
	}
	return g
}
