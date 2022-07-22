package handler

import (
	"github.com/gin-gonic/gin"
	"horizon/service"
)

func InspectionGet(c *gin.Context) {
	service.InspectionSelectByList(c)
}

func InspectionDetailGet(c *gin.Context) {
	service.InspectionSelectByInspId(c)
}
