package handler

import (
	"github.com/gin-gonic/gin"
	"horizon/service"
)

func WorkflowTemplateGet(c *gin.Context) {
	service.WorkflowTemplateSelectByList(c)
}

func WorkflowTemplatePost(c *gin.Context) {
	service.WorkflowTemplateInsert(c)
}

func WorkflowTemplatePut(c *gin.Context) {
	service.WorkflowTemplateUpdate(c)
}

func WorkflowTemplateDelete(c *gin.Context) {
	service.WorkflowTemplateDelete(c)
}

func WorkflowTemplateConfigPost(c *gin.Context) {
	service.WorkflowTemplateConfigInsert(c)
}
