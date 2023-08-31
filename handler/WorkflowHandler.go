package handler

import (
	"github.com/gin-gonic/gin"
	"horizon/service"
)

func WorkflowGet(c *gin.Context) {
	service.WorkflowSelectByList(c)
}

func WorkflowProjIdGet(c *gin.Context) {
	service.WorkflowSelectByInstId(c)
}

func WorkflowPost(c *gin.Context) {
	service.WorkflowInsert(c)
}

func WorkflowPut(c *gin.Context) {
	service.WorkflowUpdate(c)
}

func WorkflowDelete(c *gin.Context) {
	service.WorkflowDelete(c)
}
