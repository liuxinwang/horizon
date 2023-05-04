package handler

import (
	"github.com/gin-gonic/gin"
	"horizon/service"
)

func RolePermissionPost(c *gin.Context) {
	service.RolePermissionInsert(c)
}
