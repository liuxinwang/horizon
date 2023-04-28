package handler

import (
	"github.com/gin-gonic/gin"
	"horizon/service"
)

func RoleGet(c *gin.Context) {
	service.RoleSelectByList(c)
}

func RolePost(c *gin.Context) {
	service.RoleInsert(c)
}

func RolePut(c *gin.Context) {
	service.RoleUpdate(c)
}

func RoleDelete(c *gin.Context) {
	service.RoleDelete(c)
}
