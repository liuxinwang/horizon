package handler

import (
	"github.com/gin-gonic/gin"
	"horizon/service"
)

func WorkflowGet(c *gin.Context) {
	service.WorkflowSelectByList(c)
}

func WorkflowIdGet(c *gin.Context) {
	service.WorkflowSelectById(c)
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

func WorkflowIdProgressGet(c *gin.Context) {
	service.WorkflowProgressSelectById(c)
}

func WorkflowAuditPost(c *gin.Context) {
	service.WorkflowAuditUpdate(c)
}

func WorkflowCancelPost(c *gin.Context) {
	service.WorkflowCancelUpdate(c)
}

func WorkflowExecutePost(c *gin.Context) {
	service.WorkflowExecuteUpdate(c)
}

func WorkflowScheduledExecutionPost(c *gin.Context) {
	service.WorkflowScheduledExecutionUpdate(c)
}

func WorkflowIdSqlDetailGet(c *gin.Context) {
	service.WorkflowSqlDetailSelectById(c)
}
