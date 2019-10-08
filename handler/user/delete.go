package user

import (
	"github.com/gin-gonic/gin"
	"go_restful_api/handler"
	"go_restful_api/model"
	"go_restful_api/pkg/errno"
	"strconv"
)

func Delete(c *gin.Context)  {
	userId, _ := strconv.Atoi(c.Param("id"))
	if err := model.DeleteUser(uint64(userId)); err != nil {
		handler.SendResponse(c, errno.ErrDatabase, nil)
		return
	}

	handler.SendResponse(c, nil, nil)
}
