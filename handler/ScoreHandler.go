package handler

import (
	"github.com/gin-gonic/gin"
	"horizon/service"
)

func ScoreGet(c *gin.Context) {
	service.ScoreSelectByList(c)
}
