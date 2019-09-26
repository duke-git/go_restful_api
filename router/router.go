package router

import (
	"github.com/gin-gonic/gin"
	"go_restful_api/handler/sd"
	"go_restful_api/router/middleware"
	"net/http"
)
func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	//middleware
	g.Use(gin.Recovery())
	g.Use(middleware.NoCache)
	g.Use(middleware.Options)
	g.Use(middleware.Secure)
	g.Use(mw...)

	//404 handler
	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The incorrect API route")
	})

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