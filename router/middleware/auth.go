package middleware

import (
	"github.com/gin-gonic/gin"
	"go_restful_api/handler"
	"go_restful_api/pkg/errno"
	"go_restful_api/pkg/token"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, err := token.ParseRequest(c); err != nil {
			handler.SendResponse(c, errno.ErrTokenInvalid, nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
