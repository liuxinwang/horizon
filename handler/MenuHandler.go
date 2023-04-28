package handler

import (
	"github.com/gin-gonic/gin"
	"horizon/service"
)

func MenuGet(c *gin.Context) {
	service.MenuSelectByList(c)
}
