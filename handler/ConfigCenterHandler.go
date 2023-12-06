package handler

import (
	"github.com/gin-gonic/gin"
	"horizon/service"
)

func ConfigCenterInstanceGet(c *gin.Context) {
	service.ConfigCenterInstanceSelectByList(c)
}

func ConfigCenterInstanceConfigGet(c *gin.Context) {
	service.ConfigCenterInstanceConfigSelectByList(c)
}

func ConfigCenterInstanceConfigDataGet(c *gin.Context) {
	service.ConfigCenterInstanceConfigSelectByData(c)
}
