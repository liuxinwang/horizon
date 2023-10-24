package model

import (
	"database/sql"
	"fmt"
	json "github.com/json-iterator/go"
	"gorm.io/gorm"
	"horizon/config"
	"horizon/notification"
	"time"
)

type WorkflowAuditStatus string
type WorkflowStatus string
type WorkflowSqlAuditStatus string
type WorkflowSqlAuditLevel string
type WorkflowSqlExecutionStatus string

const (
	WorkflowStatusPendingAudit       WorkflowStatus = "PendingAudit"       // 待审核
	WorkflowStatusPendingExecution   WorkflowStatus = "PendingExecution"   // 待执行
	WorkflowStatusScheduledExecution WorkflowStatus = "ScheduledExecution" // 定时执行
	WorkflowStatusRejected           WorkflowStatus = "Rejected"           // 驳回
	WorkflowStatusCanceled           WorkflowStatus = "Canceled"           // 取消
	WorkflowStatusExecuting          WorkflowStatus = "Executing"          // 执行中
	WorkflowStatusExecutionFailed    WorkflowStatus = "ExecutionFailed"    // 执行失败
	WorkflowStatusFinished           WorkflowStatus = "Finished"           // 完成

	FlowAuditStatusPendingAudit  WorkflowAuditStatus = "PendingAudit" // 待审核
	FlowAuditStatusPassed        WorkflowAuditStatus = "Passed"       // 审核通过
	FlowAuditStatusAuditRejected WorkflowAuditStatus = "Rejected"     // 审核驳回

	WorkflowSqlAuditStatusPassed WorkflowSqlAuditStatus = "Passed" // 审核通过
	WorkflowSqlAuditStatusFailed WorkflowSqlAuditStatus = "Failed" // 审核失败

	WorkflowSqlAuditLevelWarning WorkflowSqlAuditLevel = "Warning" // 警告
	WorkflowSqlAuditLevelError   WorkflowSqlAuditLevel = "Error"   // 错误
	WorkflowSqlAuditLevelSuccess WorkflowSqlAuditLevel = "Success" // 成功

	WorkflowSqlExecutionStatusFailed       WorkflowSqlExecutionStatus = "Failed"       // 执行失败
	WorkflowSqlExecutionStatusSuccessfully WorkflowSqlExecutionStatus = "Successfully" // 执行成功
)

type Workflow struct {
	ID              uint             `gorm:"primaryKey;comment:主键ID" json:"id"`
	Name            string           `gorm:"type:varchar(50);not null;comment:名称" json:"name"`
	Describe        string           `gorm:"type:varchar(255);not null;comment:描述" json:"describe"`
	Status          WorkflowStatus   `gorm:"type:varchar(50);not null;default:'PendingAudit';comment:状态" json:"status"`
	ProjId          string           `gorm:"type:varchar(20);not null;comment:项目ID" json:"projId"`
	InstId          string           `gorm:"type:varchar(20);not null;comment:实例ID" json:"instId"`
	DbName          string           `gorm:"type:varchar(255);not null;comment:数据库名" json:"dbName"`
	SqlContent      string           `gorm:"type:text;not null;comment:SQL内容" json:"sqlContent"`
	UserName        string           `gorm:"type:varchar(50);not null;comment:用户名" json:"userName"`
	ScheduledAt     *time.Time       `gorm:"type:datetime;default null;comment:定时调度时间" json:"scheduledAt"`
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
	AuditStatus          WorkflowAuditStatus `gorm:"type:varchar(50);not null;default:'PendingAudit';comment:状态" json:"auditStatus"`
	IsAudit              uint                `gorm:"not null;default:0;comment:审核标识（0：未审核，1：已审核）" json:"isAudit"`
	CreatedAt            time.Time           `gorm:"type:datetime;not null;default:current_timestamp;comment:创建时间" json:"createdAt"`
	UpdatedAt            time.Time           `gorm:"type:datetime;not null;default:current_timestamp on update current_timestamp;comment:修改时间" json:"updatedAt"`
	User                 User                `gorm:"foreignKey:AssigneeUserName;references:UserName" json:"user"`
}

type WorkflowSqlDetail struct {
	ID              uint                       `gorm:"primaryKey;comment:主键ID" json:"id"`
	WorkflowId      uint                       `gorm:"not null;comment:工单ID" json:"workflowId"`
	SerialNumber    uint                       `gorm:"not null;comment:工单语句序号" json:"serialNumber"`
	Statement       string                     `gorm:"type:text;not null;comment:工单语句" json:"statement"`
	AuditStatus     WorkflowSqlAuditStatus     `gorm:"type:varchar(20);not null;comment:审核状态" json:"auditStatus"`
	AuditLevel      WorkflowSqlAuditLevel      `gorm:"type:varchar(20);not null;comment:审核等级" json:"auditLevel"`
	AuditMsg        string                     `gorm:"type:varchar(1000);not null;comment:审核信息" json:"auditMsg"`
	ExecutionStatus WorkflowSqlExecutionStatus `gorm:"type:varchar(20);not null;comment:执行状态" json:"executionStatus"`
	ExecutionMsg    string                     `gorm:"type:varchar(1000);not null;comment:执行信息" json:"executionMsg"`
	CreatedAt       time.Time                  `gorm:"type:datetime;not null;default:current_timestamp;comment:创建时间" json:"createdAt"`
	UpdatedAt       time.Time                  `gorm:"type:datetime;not null;default:current_timestamp on update current_timestamp;comment:修改时间" json:"updatedAt"`
}

func (w *Workflow) AfterUpdate(tx *gorm.DB) (err error) {
	var workflowUser User
	result := tx.Where("user_name = ?", w.UserName).First(&workflowUser)
	if result.Error != nil {
		return result.Error
	}

	if w.Status == WorkflowStatusPendingExecution {
		markdown := notification.DingContentMarkdown{
			Title: "SQL工单通知",
			Text: fmt.Sprintf(
				"用户：@%s 提交的 SQL工单【%s】，已审核通过。[工单详情](%s/sqlaudit/workflowDetail/%d)",
				workflowUser.Phone, w.Name, config.Conf.General.HomeAddress, w.ID),
		}
		at := notification.DingContentAt{
			AtMobiles: []string{workflowUser.Phone},
			AtUserIds: []string{},
			IsAtAll:   false,
		}
		content := notification.DingContent{
			MsgType:         "markdown",
			ContentMarkdown: markdown,
			ContentAt:       at,
		}
		marshal, err := json.Marshal(content)
		if err != nil {
			return err
		}
		notification.SendDingDing(string(marshal))
	}
	if w.Status == WorkflowStatusRejected {
		markdown := notification.DingContentMarkdown{
			Title: "SQL工单通知",
			Text: fmt.Sprintf(
				"用户：@%s 提交的 SQL工单【%s】，被驳回！[工单详情](%s/sqlaudit/workflowDetail/%d)",
				workflowUser.Phone, w.Name, config.Conf.General.HomeAddress, w.ID),
		}
		at := notification.DingContentAt{
			AtMobiles: []string{workflowUser.Phone},
			AtUserIds: []string{},
			IsAtAll:   false,
		}
		content := notification.DingContent{
			MsgType:         "markdown",
			ContentMarkdown: markdown,
			ContentAt:       at,
		}
		marshal, err := json.Marshal(content)
		if err != nil {
			return err
		}
		notification.SendDingDing(string(marshal))
	}
	if w.Status == WorkflowStatusFinished {
		markdown := notification.DingContentMarkdown{
			Title: "SQL工单通知",
			Text: fmt.Sprintf(
				"用户：@%s 提交的 SQL工单【%s】，已执行成功。[工单详情](%s/sqlaudit/workflowDetail/%d)",
				workflowUser.Phone, w.Name, config.Conf.General.HomeAddress, w.ID),
		}
		at := notification.DingContentAt{
			AtMobiles: []string{workflowUser.Phone},
			AtUserIds: []string{},
			IsAtAll:   false,
		}
		content := notification.DingContent{
			MsgType:         "markdown",
			ContentMarkdown: markdown,
			ContentAt:       at,
		}
		marshal, err := json.Marshal(content)
		if err != nil {
			return err
		}
		notification.SendDingDing(string(marshal))
	}
	if w.Status == WorkflowStatusExecutionFailed {
		markdown := notification.DingContentMarkdown{
			Title: "SQL工单通知",
			Text: fmt.Sprintf(
				"用户：@%s 提交的 SQL工单【%s】，执行失败！[工单详情](%s/sqlaudit/workflowDetail/%d)",
				workflowUser.Phone, w.Name, config.Conf.General.HomeAddress, w.ID),
		}
		at := notification.DingContentAt{
			AtMobiles: []string{workflowUser.Phone},
			AtUserIds: []string{},
			IsAtAll:   false,
		}
		content := notification.DingContent{
			MsgType:         "markdown",
			ContentMarkdown: markdown,
			ContentAt:       at,
		}
		marshal, err := json.Marshal(content)
		if err != nil {
			return err
		}
		notification.SendDingDing(string(marshal))
	}
	return
}

func (wr *WorkflowRecord) AfterCreate(tx *gorm.DB) (err error) {
	var workflow Workflow
	result := tx.Where("id = ?", wr.WorkflowId).First(&workflow)
	if result.Error != nil {
		return result.Error
	}
	var workflowUser User
	result = tx.Where("user_name = ?", workflow.UserName).First(&workflowUser)
	if result.Error != nil {
		return result.Error
	}

	var workflowAssigneeUser User
	result = tx.Where("user_name = ?", wr.AssigneeUserName).First(&workflowAssigneeUser)
	if result.Error != nil {
		return result.Error
	}

	if wr.AuditStatus == FlowAuditStatusPendingAudit {
		markdown := notification.DingContentMarkdown{
			Title: "SQL工单通知",
			Text: fmt.Sprintf(
				"用户：%s 提交的 SQL工单【%s】，正在等待 @%s 审批，请确认！[审核](%s/sqlaudit/workflowDetail/%d)",
				workflowUser.NickName, workflow.Name,
				workflowAssigneeUser.Phone, config.Conf.General.HomeAddress, wr.WorkflowId),
		}
		at := notification.DingContentAt{
			AtMobiles: []string{workflowAssigneeUser.Phone},
			AtUserIds: []string{},
			IsAtAll:   false,
		}
		content := notification.DingContent{
			MsgType:         "markdown",
			ContentMarkdown: markdown,
			ContentAt:       at,
		}
		marshal, err := json.Marshal(content)
		if err != nil {
			return err
		}
		notification.SendDingDing(string(marshal))
	}
	return nil
}
