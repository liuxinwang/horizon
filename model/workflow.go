package model

import "time"

type workflowStatus string

const (
	WorkflowStatusWaitForAudit     workflowStatus = "wait_for_audit"
	WorkflowStatusWaitForExecution workflowStatus = "wait_for_execution"
	WorkflowStatusReject           workflowStatus = "rejected"
	WorkflowStatusCancel           workflowStatus = "canceled"
	WorkflowStatusExecuting        workflowStatus = "executing"
	WorkflowStatusExecFailed       workflowStatus = "exec_failed"
	WorkflowStatusFinish           workflowStatus = "finished"
)

type Workflow struct {
	ID        uint           `gorm:"primaryKey;comment:主键ID" json:"id"`
	Name      string         `gorm:"type:varchar(50);not null;comment:名称" json:"name"`
	Describe  string         `gorm:"type:varchar(255);not null;comment:描述" json:"describe"`
	Status    workflowStatus `gorm:"type:varchar(50);not null;default:'wait_for_audit';comment:状态" json:"status"`
	ProjId    string         `gorm:"type:varchar(20);not null;comment:项目ID" json:"projId"`
	UserName  string         `gorm:"type:varchar(50);not null;comment:用户名" json:"userName"`
	CreatedAt time.Time      `gorm:"type:datetime;not null;default:current_timestamp;comment:创建时间" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"type:datetime;not null;default:current_timestamp on update current_timestamp;comment:修改时间" json:"updatedAt"`
}
