package handler

import (
	"github.com/gin-gonic/gin"
	"horizon/service"
)

func RolePermissionGet(c *gin.Context) {
	service.RolePermissionSelectByList(c)
}
