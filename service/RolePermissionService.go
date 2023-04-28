package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// RolePermissionSelectByList 查询角色权限列表
func RolePermissionSelectByList(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": "", "err": ""})
}
