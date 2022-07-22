package handler

import (
	"github.com/gin-gonic/gin"
	"horizon/service"
)

func InstanceGet(c *gin.Context) {
	service.InstanceSelectByList(c)
}

func InstancePost(c *gin.Context) {
	service.InstanceInsert(c)
}

func InstancePut(c *gin.Context) {
	service.InstanceUpdate(c)
}

func InstanceDelete(c *gin.Context) {
	service.InstanceDelete(c)
}
