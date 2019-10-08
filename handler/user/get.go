package user

import (
	"github.com/gin-gonic/gin"
	"go_restful_api/handler"
	"go_restful_api/model"
	"go_restful_api/pkg/errno"
)

func Get(c *gin.Context) {
	username := c.Param("username")
	// Get the user by the `username` from the database.
	user, err := model.GetUser(username)
	if err != nil {
		handler.SendResponse(c, errno.ErrUserNotFound, nil)
		return
	}

	handler.SendResponse(c, nil, user)
}