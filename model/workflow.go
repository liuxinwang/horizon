package model

import (
	"database/sql"
	"time"
)

type workflowAuditStatus string
type workflowStatus string
type WorkflowSqlAuditStatus string
type WorkflowSqlAuditLevel string

const (
	WorkflowStatusPendingAudit     workflowStatus = "PendingAudit"
	WorkflowStatusPendingExecution workflowStatus = "PendingExecution"
	WorkflowStatusRejected         workflowStatus = "Rejected"
	WorkflowStatusCanceled         workflowStatus = "Canceled"
	WorkflowStatusExecuting        workflowStatus = "Executing"
	WorkflowStatusExecutionFailed  workflowStatus = "ExecutionFailed"
	WorkflowStatusFinished         workflowStatus = "Finished"

	FlowAuditStatusPendingAudit  workflowAuditStatus = "PendingAudit"
	FlowAuditStatusPassed        workflowAuditStatus = "Passed"
	FlowAuditStatusAuditRejected workflowAuditStatus = "Rejected"

	WorkflowSqlAuditStatusPassed WorkflowSqlAuditStatus = "Passed"
	WorkflowSqlAuditStatusFailed WorkflowSqlAuditStatus = "Failed"

	WorkflowSqlAuditLevelWarning WorkflowSqlAuditLevel = "Warning"
	WorkflowSqlAuditLevelError   WorkflowSqlAuditLevel = "Error"
	WorkflowSqlAuditLevelSuccess WorkflowSqlAuditLevel = "Success"
)

type Workflow struct {
	ID              uint             `gorm:"primaryKey;comment:主键ID" json:"id"`
	Name            string           `gorm:"type:varchar(50);not null;comment:名称" json:"name"`
	Describe        string           `gorm:"type:varchar(255);not null;comment:描述" json:"describe"`
	Status          workflowStatus   `gorm:"type:varchar(50);not null;default:'PendingAudit';comment:状态" json:"status"`
	ProjId          string           `gorm:"type:varchar(20);not null;comment:项目ID" json:"projId"`
	InstId          string           `gorm:"type:varchar(20);not null;comment:实例ID" json:"instId"`
	DbName          string           `gorm:"type:varchar(255);not null;comment:数据库名" json:"dbName"`
	SqlContent      string           `gorm:"type:text;not null;comment:SQL内容" json:"sqlContent"`
	UserName        string           `gorm:"type:varchar(50);not null;comment:用户名" json:"userName"`
	CreatedAt       time.Time        `gorm:"type:datetime;not null;default:current_timestamp;comment:创建时间" json:"createdAt"`
	UpdatedAt       time.Time        `gorm:"type:datetime;not null;default:current_timestamp on update current_timestamp;comment:修改时间" json:"updatedAt"`
	WorkflowRecords []WorkflowRecord `gorm:"foreignKey:WorkflowId;references:ID" json:"workflowRecords"`
}

type WorkflowTemplate struct {
	ID                      uint                     `gorm:"primaryKey;comment:主键ID" json:"id"`
	Code                    uint                     `gorm:"uniqueIndex:uniq_code;not null;comment:编号" json:"code"`
	Name                    string                   `gorm:"type:varchar(50);not null;comment:名称" json:"name"`
	CreatedAt               time.Time                `gorm:"type:datetime;not null;default:current_timestamp;comment:创建时间" json:"createdAt"`
	UpdatedAt               time.Time                `gorm:"type:datetime;not null;default:current_timestamp on update current_timestamp;comment:修改时间" json:"updatedAt"`
	WorkflowTemplateDetails []WorkflowTemplateDetail `gorm:"foreignKey:WorkflowTemplateId;references:ID" json:"workflowTemplateDetails"`
	// Projects                []Project                `gorm:"foreignKey:WorkflowTemplateCode;references:Code" json:"projects"`
}

type WorkflowTemplateDetail struct {
	ID                   uint      `gorm:"primaryKey;comment:主键ID" json:"id"`
	WorkflowTemplateId   uint      `gorm:"not null;comment:工作流模版ID" json:"workflowTemplateId"`
	WorkflowTemplateCode uint      `gorm:"not null;comment:工作流模版编号" json:"workflowTemplateCode"`
	SerialNumber         uint      `gorm:"not null;comment:工作流序号" json:"serialNumber"`
	NodeName             string    `gorm:"type:varchar(20);not null;comment:节点名称" json:"nodeName"`
	ProjectRoleId        string    `gorm:"type:varchar(50);not null;comment:项目角色ID" json:"projectRoleId"`
	CreatedAt            time.Time `gorm:"type:datetime;not null;default:current_timestamp;comment:创建时间" json:"createdAt"`
	UpdatedAt            time.Time `gorm:"type:datetime;not null;default:current_timestamp on update current_timestamp;comment:修改时间" json:"updatedAt"`
}

type WorkflowRecord struct {
	ID                   uint                `gorm:"primaryKey;comment:主键ID" json:"id"`
	WorkflowId           uint                `gorm:"not null;comment:工单ID" json:"workflowId"`
	WorkflowTemplateCode uint                `gorm:"not null;comment:工作流模版Code" json:"workflowTemplateCode"`
	FlowNodeName         string              `gorm:"type:varchar(20);not null;comment:节点名称" json:"flowNodeName"`
	FlowSerialNumber     uint                `gorm:"not null;comment:工作流序号" json:"flowSerialNumber"`
	AssigneeUserName     string              `gorm:"type:varchar(50);not null;comment:受理用户" json:"assigneeUserName"`
	HandledAt            sql.NullTime        `gorm:"type:datetime;default null;comment:处理时间" json:"handledAt"`
	Remarks              string              `gorm:"type:varchar(255);not null;comment:处理结果/备注" json:"remarks"`
	AuditStatus          workflowAuditStatus `gorm:"type:varchar(50);not null;default:'PendingAudit';comment:状态" json:"auditStatus"`
	IsAudit              uint                `gorm:"not null;default:0;comment:审核标识（0：未审核，1：已审核）" json:"isAudit"`
	CreatedAt            time.Time           `gorm:"type:datetime;not null;default:current_timestamp;comment:创建时间" json:"createdAt"`
	UpdatedAt            time.Time           `gorm:"type:datetime;not null;default:current_timestamp on update current_timestamp;comment:修改时间" json:"updatedAt"`
}

type WorkflowSqlDetail struct {
	ID           uint                   `gorm:"primaryKey;comment:主键ID" json:"id"`
	WorkflowId   uint                   `gorm:"not null;comment:工单ID" json:"workflowId"`
	SerialNumber uint                   `gorm:"not null;comment:工单语句序号" json:"serialNumber"`
	Statement    string                 `gorm:"type:text;not null;comment:工单语句" json:"statement"`
	AuditStatus  WorkflowSqlAuditStatus `gorm:"type:varchar(20);not null;comment:审核状态" json:"auditStatus"`
	AuditLevel   WorkflowSqlAuditLevel  `gorm:"type:varchar(20);not null;comment:审核等级" json:"auditLevel"`
	AuditMsg     string                 `gorm:"type:varchar(1000);not null;comment:审核信息" json:"auditMsg"`
	CreatedAt    time.Time              `gorm:"type:datetime;not null;default:current_timestamp;comment:创建时间" json:"createdAt"`
	UpdatedAt    time.Time              `gorm:"type:datetime;not null;default:current_timestamp on update current_timestamp;comment:修改时间" json:"updatedAt"`
}
