package service

import (
	"github.com/gin-gonic/gin"
	_ "github.com/pingcap/tidb/types/parser_driver"
)

type Workflow interface {
	WorkflowSelectByList(c *gin.Context)             // 查询列表
	WorkflowSelectById(c *gin.Context)               // 查看工单信息
	WorkflowInsert(c *gin.Context)                   // 新增
	WorkflowUpdate(c *gin.Context)                   // 修改
	WorkflowDelete(c *gin.Context)                   // 删除
	WorkflowProgressSelectById(c *gin.Context)       //  查询工作流进度
	WorkflowAuditUpdate(c *gin.Context)              // 工单审核
	WorkflowCancelUpdate(c *gin.Context)             // 工单取消
	WorkflowExecuteUpdate(c *gin.Context)            // 工单执行
	WorkflowScheduledExecutionUpdate(c *gin.Context) // 工单定时执行
	WorkflowSqlDetailSelectById(c *gin.Context)      // 查询工单sql审核明细
}
