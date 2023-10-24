package service

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/pingcap/tidb/types/parser_driver"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"horizon/config"
	"horizon/model"
	"horizon/plugins"
	"horizon/utils"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// WorkflowSelectByList 查询列表
func WorkflowSelectByList(c *gin.Context) {
	// 变量初始化
	pageNo, _ := strconv.Atoi(c.Query("pageNo"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	var Db = model.Db
	type Workflow struct {
		model.Workflow
		InstName string `json:"instName"`
		NickName string `json:"nickName"`
	}
	var workflows []Workflow
	var totalCount int64
	data := gin.H{"totalCount": 0, "data": &[]model.Workflow{}, "pageNo": pageNo, "pageSize": pageSize, "totalPage": 0}

	// 查询条件处理
	if name, isExist := c.GetQuery("name"); isExist == true && strings.Trim(name, " ") != "" {
		Db = Db.Where("workflows.name like ?", "%"+name+"%")
	}
	if status, isExist := c.GetQuery("status"); isExist == true && status != "0" {
		Db = Db.Where("workflows.status = ?", status)
	}

	userInfo, _ := c.Keys["UserName"]
	userName := userInfo.(*model.User).UserName
	var user model.User
	model.Db.Preload("UserRoles").Preload("ProjectUsers").Where("user_name = ?", userName).First(&user)
	// 系统DBA、系统admin 可以看到所有工单数据

	// 项目负责人、项目DBA 可以看到所属项目的所有工单数据
	// 用户默认查询用户自己的工单数据
	// 用户还可以查询参与审核的工单

	// 执行查询
	Db = Db.Preload("WorkflowRecords", func(db *gorm.DB) *gorm.DB {
		return db.Order("workflow_records.id asc")
	}).Preload("WorkflowRecords.User")
	Db.Select("workflows.*, instances.name as inst_name, users.nick_name").
		Joins("left join instances on workflows.inst_id = instances.inst_id").
		Joins("left join users on workflows.user_name = users.user_name").
		Where("(workflows.user_name = ?", userName).
		Or("workflows.id in (?))", model.Db.Table("workflow_records").
			Where("assignee_user_name = ?", userName).Select("workflow_id")).
		Order("created_at desc").
		Limit(pageSize).Offset((pageNo-1)*pageSize - 1).
		Find(&workflows).
		Limit(-1).
		Offset(-1).
		Count(&totalCount)

	// 处理结果集并返回
	totalPage := math.Ceil(float64(totalCount) / float64(pageSize))
	if totalCount > 0 {
		data = gin.H{"totalCount": totalCount, "data": &workflows, "pageNo": pageNo, "pageSize": pageSize, "totalPage": totalPage}
	}
	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": data, "err": ""})
}

// WorkflowSelectById 查看工单信息
func WorkflowSelectById(c *gin.Context) {
	var Db = model.Db
	var workflow model.Workflow
	workflowId, _ := strconv.Atoi(c.Param("id"))

	userInfo, _ := c.Keys["UserName"]
	userName := userInfo.(*model.User).UserName

	canQueryDetail := CanQueryWorkflowDetail(workflowId, userName)
	if !canQueryDetail {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": "无工单详情查询权限！"})
		return
	}

	// 执行查询
	Db.Where("id = ?", workflowId).Find(&workflow)
	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": &workflow, "err": ""})
}

// WorkflowInsert 新增
func WorkflowInsert(c *gin.Context) {
	// 参数映射到对象
	var workflow model.Workflow
	if err := c.ShouldBind(&workflow); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": err.Error()})
		return
	}
	userInfo, _ := c.Keys["UserName"]
	workflow.UserName = userInfo.(*model.User).UserName

	// 判断是否有新增工单权限
	canAdd := CanAddWorkflow(&workflow)
	if !canAdd {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": "无项目新增工单权限！"})
		return
	}

	tx := model.Db.Begin()

	var instance model.Instance
	tx.Where("inst_id = ?", workflow.InstId).First(&instance)

	// insert workflow
	result := tx.Create(&workflow)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": result.Error.Error()})
		tx.Rollback()
		return
	}

	var plugin plugins.Plugin
	switch instance.Type {
	case model.MySQLType:
		plugin = &plugins.MySQLPlugin{Instance: &instance}
	case model.DorisType:
		plugin = &plugins.DorisPlugin{Instance: &instance}
	default:
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": "不支持的数据库类型：" + instance.Type})
		tx.Rollback()
		return
	}

	auditResult, err := plugin.Audit(&workflow)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": err.Error()})
		tx.Rollback()
		return
	}
	for i, ar := range auditResult {
		var workflowSqlDetail model.WorkflowSqlDetail
		workflowSqlDetail.WorkflowId = workflow.ID
		workflowSqlDetail.SerialNumber = uint(i + 1)
		workflowSqlDetail.Statement = ar.Statement
		workflowSqlDetail.AuditStatus = ar.AuditStatus
		workflowSqlDetail.AuditLevel = ar.AuditLevel
		workflowSqlDetail.AuditMsg = ar.AuditMsg
		createRs := tx.Create(&workflowSqlDetail)
		if createRs.Error != nil {
			c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": err.Error()})
			tx.Rollback()
			return
		}
		if ar.AuditStatus != model.WorkflowSqlAuditStatusPassed ||
			ar.AuditLevel == model.WorkflowSqlAuditLevelError {
			c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": workflowSqlDetail.AuditMsg})
			tx.Rollback()
			return
		}
	}

	// insert WorkflowRecord
	var workflowRecord model.WorkflowRecord
	workflowRecord.WorkflowId = workflow.ID
	// 获取WorkflowTemplateCode
	var project model.Project
	tx.Where("proj_id = ?", workflow.ProjId).First(&project)
	workflowRecord.WorkflowTemplateCode = project.WorkflowTemplateCode
	var workflowTemplate model.WorkflowTemplate
	tx.Where("code = ?", project.WorkflowTemplateCode).First(&workflowTemplate)
	// 获取first nodeName
	var workflowTemplateDetail model.WorkflowTemplateDetail
	tx.Where("workflow_template_id = ?", workflowTemplate.ID).Order("serial_number").First(&workflowTemplateDetail)
	workflowRecord.FlowNodeName = workflowTemplateDetail.NodeName
	workflowRecord.FlowSerialNumber = workflowTemplateDetail.SerialNumber
	// 获取受理用户
	var projectUser model.ProjectUser
	tx.Where("role_id = ?", workflowTemplateDetail.ProjectRoleId).First(&projectUser)
	workflowRecord.AssigneeUserName = projectUser.UserName
	result = tx.Create(&workflowRecord)
	tx.Commit()
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": result.Error.Error()})
		tx.Rollback()
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": "", "err": ""})
	}
}

// WorkflowUpdate 重新修改
func WorkflowUpdate(c *gin.Context) {
	WorkflowInsert(c)
}

// WorkflowDelete 删除
func WorkflowDelete(c *gin.Context) {
	id := c.Param("id")
	result := model.Db.Delete(&model.Workflow{}, id)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": result.Error.Error()})
	} else if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": "工单不存在"})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": "", "err": ""})
	}
}

// WorkflowProgressSelectById 查询工作流进度
func WorkflowProgressSelectById(c *gin.Context) {
	type Result struct {
		SerialNumber         uint         `json:"serialNumber"`
		NodeName             string       `json:"nodeName"`
		ProjRoleName         string       `json:"projRoleName"`
		WorkflowTemplateCode uint         `json:"workflowTemplateCode"`
		RecordId             uint         `json:"recordId"`
		AssigneeUserName     string       `json:"assigneeUserName"`
		HandledAt            sql.NullTime `json:"handledAt"`
		Remarks              string       `json:"remarks"`
		AuditStatus          string       `json:"auditStatus"`
		IsAudit              uint         `json:"isAudit"`
		CreatedAt            time.Time    `json:"createdAt"`
	}
	var results []Result
	progressSql := "select \n" +
		"wtd.serial_number, wtd.node_name, pr.name as proj_role_name, p.workflow_template_code, \n" +
		"wr.id as record_id, wr.assignee_user_name, wr.handled_at, wr.remarks, wr.audit_status, wr.is_audit, wr.created_at\n" +
		"from workflows w\n" +
		"inner join projects p on w.proj_id = p.proj_id\n" +
		"inner join workflow_templates wt on p.workflow_template_code = wt.code\n" +
		"inner join workflow_template_details wtd on wt.id = wtd.workflow_template_id\n" +
		"inner join project_roles pr on wtd.project_role_id = pr.id\n" +
		"left join workflow_records wr on w.id = wr.workflow_id and wr.flow_serial_number = wtd.serial_number\n" +
		"where w.id = ?\n" +
		"order by wtd.serial_number asc;"
	result := model.Db.Raw(progressSql, c.Param("id")).Scan(&results)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": &results, "err": ""})
}

// WorkflowAuditUpdate 工单审核
func WorkflowAuditUpdate(c *gin.Context) {
	// 参数映射到对象
	var workflowRecord model.WorkflowRecord
	if err := c.ShouldBind(&workflowRecord); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": err.Error()})
		return
	}

	userInfo, _ := c.Keys["UserName"]
	userName := userInfo.(*model.User).UserName

	// 判断是否有审核权限
	canAudit := CanAuditWorkflow(&workflowRecord, userName)
	if !canAudit {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": "无工单审核权限！"})
		return
	}

	// 更新workflowRecord状态
	// auditStatus = FlowAuditStatusPassed
	// 判断当前审批节点
	// 是最后一个，更新workflow状态为审核完成，待上线
	// 不是最后一个，获取下一个审批节点，插入workflowRecord
	tx := model.Db.Begin()
	var workflow model.Workflow
	tx.Where("id = ?", workflowRecord.WorkflowId).First(&workflow)
	if workflow.Status != model.WorkflowStatusPendingAudit {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": "非待审核工单！"})
		return
	}
	switch workflowRecord.AuditStatus {
	case model.FlowAuditStatusPassed:
		var workflowTemplateDetail model.WorkflowTemplateDetail
		tx.Where("workflow_template_code = ? and  serial_number > ?", workflowRecord.WorkflowTemplateCode, workflowRecord.FlowSerialNumber).
			Order("serial_number asc").First(&workflowTemplateDetail)
		if workflowTemplateDetail.ID > 0 {
			// 不是最后一个
			tx.Model(&model.WorkflowRecord{}).Where("id = ?", workflowRecord.ID).
				Updates(model.WorkflowRecord{
					AuditStatus: model.FlowAuditStatusPassed,
					IsAudit:     1,
					Remarks:     workflowRecord.Remarks,
					HandledAt:   sql.NullTime{Time: time.Now(), Valid: true}})
			// insert WorkflowRecord
			var nextWorkflowRecord model.WorkflowRecord
			nextWorkflowRecord.WorkflowId = workflowRecord.WorkflowId
			nextWorkflowRecord.WorkflowTemplateCode = workflowRecord.WorkflowTemplateCode
			nextWorkflowRecord.FlowNodeName = workflowTemplateDetail.NodeName
			nextWorkflowRecord.FlowSerialNumber = workflowTemplateDetail.SerialNumber
			// 获取受理用户
			var projectUser model.ProjectUser
			tx.Where("role_id = ?", workflowTemplateDetail.ProjectRoleId).First(&projectUser)
			nextWorkflowRecord.AssigneeUserName = projectUser.UserName
			result := tx.Create(&nextWorkflowRecord)
			if result.Error != nil {
				c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": result.Error})
				tx.Rollback()
				return
			}
		} else {
			// 是最后一个
			tx.Model(&model.WorkflowRecord{}).Where("id = ?", workflowRecord.ID).
				Updates(model.WorkflowRecord{
					AuditStatus: model.FlowAuditStatusPassed,
					IsAudit:     1,
					Remarks:     workflowRecord.Remarks,
					HandledAt:   sql.NullTime{Time: time.Now(), Valid: true}})
			result := tx.Model(&workflow).Update("status", model.WorkflowStatusPendingExecution)
			if result.Error != nil {
				c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": result.Error})
				tx.Rollback()
				return
			}
		}
	case model.FlowAuditStatusAuditRejected:
		// 更新workflowRecord状态
		// auditStatus = FlowAuditStatusAuditRejected
		result := tx.Model(&model.WorkflowRecord{}).Where("id = ?", workflowRecord.ID).
			Updates(model.WorkflowRecord{
				AuditStatus: model.FlowAuditStatusAuditRejected,
				IsAudit:     1,
				Remarks:     workflowRecord.Remarks,
				HandledAt:   sql.NullTime{Time: time.Now(), Valid: true}})
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": result.Error.Error()})
			tx.Rollback()
			return
		}
		result = tx.Model(&workflow).Where("id = ?", workflowRecord.WorkflowId).
			Updates(model.Workflow{Status: model.WorkflowStatusRejected})
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": result.Error.Error()})
			tx.Rollback()
			return
		}
	default:
	}
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": "", "err": ""})
}

// WorkflowCancelUpdate 工单取消
func WorkflowCancelUpdate(c *gin.Context) {
	// 参数映射到对象
	var workflow model.Workflow
	if err := c.ShouldBind(&workflow); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": err.Error()})
		return
	}

	userInfo, _ := c.Keys["UserName"]
	userName := userInfo.(*model.User).UserName

	canCancel := CanCancelWorkflow(&workflow, userName)
	if !canCancel {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": "无工单取消权限！"})
		return
	}
	var workflowStatus string
	model.Db.Select("status").Model(&workflow).Where("id = ?", workflow.ID).First(&workflowStatus)
	if model.WorkflowStatus(workflowStatus) == model.WorkflowStatusRejected ||
		model.WorkflowStatus(workflowStatus) == model.WorkflowStatusCanceled ||
		model.WorkflowStatus(workflowStatus) == model.WorkflowStatusExecuting ||
		model.WorkflowStatus(workflowStatus) == model.WorkflowStatusExecutionFailed ||
		model.WorkflowStatus(workflowStatus) == model.WorkflowStatusFinished {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": "工单状态不支持取消！"})
		return
	}

	result := model.Db.Model(&workflow).Where("id = ?", workflow.ID).Updates(&model.Workflow{Status: model.WorkflowStatusCanceled})
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": "", "err": ""})
}

// WorkflowExecuteUpdate 工单执行
func WorkflowExecuteUpdate(c *gin.Context) {
	var workflow model.Workflow
	if err := c.ShouldBind(&workflow); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": err.Error()})
		return
	}

	userInfo, _ := c.Keys["UserName"]
	userName := userInfo.(*model.User).UserName

	canExecute := CanExecuteWorkflow(&workflow, userName)
	if !canExecute {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": "非提交人，无工单执行权限！"})
		return
	}

	var workflowStatus string
	model.Db.Select("status").Model(&workflow).Where("id = ?", workflow.ID).First(&workflowStatus)
	if model.WorkflowStatus(workflowStatus) != model.WorkflowStatusPendingExecution {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": "工单状态不支持执行！"})
		return
	}

	// 获取实例信息
	var instance model.Instance
	result := model.Db.First(&instance, "inst_id = ?", workflow.InstId)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": fmt.Sprintf("%v 实例不存在", workflow.InstId)})
		return
	}
	// 更新状态 WorkflowStatusExecuting
	result = model.Db.Model(&workflow).Where("id = ?", workflow.ID).Updates(&model.Workflow{Status: model.WorkflowStatusExecuting})
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": result.Error.Error()})
		return
	}

	var plugin plugins.Plugin
	switch instance.Type {
	case model.MySQLType:
		plugin = &plugins.MySQLPlugin{Instance: &instance}
	case model.DorisType:
		plugin = &plugins.DorisPlugin{Instance: &instance}
	}

	go func() {
		err := plugin.Execute(&workflow)
		if err != nil {
			// 更新状态 WorkflowStatusExecutionFailed
			result = model.Db.Model(&workflow).Where("id = ?", workflow.ID).Updates(&model.Workflow{Status: model.WorkflowStatusExecutionFailed})
			if result.Error != nil {
				// c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": result.Error.Error()})
				return
			}
			// c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": err.Error()})
			return
		}

		// 更新状态 WorkflowStatusFinished
		result = model.Db.Model(&workflow).Where("id = ?", workflow.ID).Updates(&model.Workflow{Status: model.WorkflowStatusFinished})
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": result.Error.Error()})
			return
		}
	}()
	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": "", "err": ""})
}

// WorkflowScheduledExecutionUpdate 工单定时执行
func WorkflowScheduledExecutionUpdate(c *gin.Context) {
	// 参数映射到对象
	var workflow model.Workflow
	if err := c.ShouldBind(&workflow); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": err.Error()})
		return
	}

	userInfo, _ := c.Keys["UserName"]
	userName := userInfo.(*model.User).UserName

	canExecute := CanExecuteWorkflow(&workflow, userName)
	if !canExecute {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": "非提交人，无工单执行权限！"})
		return
	}

	var workflowStatus string
	model.Db.Select("status").Model(&workflow).Where("id = ?", workflow.ID).First(&workflowStatus)
	if model.WorkflowStatus(workflowStatus) != model.WorkflowStatusPendingExecution {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": "工单状态不支持执行！"})
		return
	}

	result := model.Db.Model(&workflow).Where("id = ?", workflow.ID).Updates(&model.Workflow{
		Status:      model.WorkflowStatusScheduledExecution,
		ScheduledAt: workflow.ScheduledAt,
	})
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": "", "err": ""})
}

// WorkflowSqlDetailSelectById 查询工单sql审核明细
func WorkflowSqlDetailSelectById(c *gin.Context) {
	// 变量初始化
	pageNo, _ := strconv.Atoi(c.Query("pageNo"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	workflowId, _ := strconv.Atoi(c.Query("workflowId"))

	userInfo, _ := c.Keys["UserName"]
	userName := userInfo.(*model.User).UserName

	canQueryDetail := CanQueryWorkflowDetail(workflowId, userName)
	if !canQueryDetail {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": "无工单详情查询权限！"})
		return
	}

	var workflowSqlDetails []model.WorkflowSqlDetail
	var totalCount int64
	data := gin.H{"totalCount": 0, "data": &[]model.WorkflowSqlDetail{}, "pageNo": pageNo, "pageSize": pageSize, "totalPage": 0}
	var Db = model.Db
	result := Db.Where("workflow_id = ?", workflowId).Find(&workflowSqlDetails)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": result.Error.Error()})
		return
	}
	Db.Model(&model.WorkflowSqlDetail{}).Where("workflow_id = ?", workflowId).Count(&totalCount)
	// 处理结果集并返回
	totalPage := math.Ceil(float64(totalCount) / float64(pageSize))
	if totalCount > 0 {
		data = gin.H{"totalCount": totalCount, "data": &workflowSqlDetails, "pageNo": pageNo, "pageSize": pageSize, "totalPage": totalPage}
	}

	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": data, "err": ""})
}

// ExecuteSQL 执行实例SQL
func ExecuteSQL(instance *model.Instance, db string, sql string) error {
	dsn := "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=1s"
	dsn = fmt.Sprintf(dsn, instance.User, utils.DecryptAES([]byte(config.Conf.General.SecretKey), instance.Password), instance.Ip, instance.Port, db)
	Db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	result := Db.Exec(sql)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// CanAddWorkflow 判断用户当前是否可新增
func CanAddWorkflow(workflow *model.Workflow) bool {
	// 工单项目 in 用户所属项目
	var count int64
	model.Db.Model(&model.ProjectUser{}).
		Where("user_name = ? and proj_id = ?", workflow.UserName, workflow.ProjId).Count(&count)
	if count > 0 {
		return true
	}
	return false
}

// CanAuditWorkflow 判断用户当前是否可审核
func CanAuditWorkflow(workflowRecord *model.WorkflowRecord, userName string) bool {
	var assigneeUserName string
	model.Db.Model(&model.WorkflowRecord{}).Select("assignee_user_name").
		Where("id = ?", workflowRecord.ID).First(&assigneeUserName)
	// 用户属于当前工单审批人
	if assigneeUserName == userName {
		return true
	}
	return false
}

// CanCancelWorkflow 判断用户当前是否可取消
func CanCancelWorkflow(workflow *model.Workflow, userName string) bool {
	// 工单提交人 = 登录人
	if workflow.UserName == userName {
		return true
	}
	return false
}

// CanExecuteWorkflow 判断用户当前是否可执行
func CanExecuteWorkflow(workflow *model.Workflow, userName string) bool {
	// 工单提交人 = 登录人
	if workflow.UserName == userName {
		return true
	}
	return false
}

// CanQueryWorkflowDetail 判断用户当前是否可查询工单sql详情
func CanQueryWorkflowDetail(workflowId int, userName string) bool {
	// 登录人 in （工单提交人，工单审核人）
	var workflow *model.Workflow
	model.Db.Model(&model.Workflow{}).Where("id = ?", workflowId).First(&workflow)
	if workflow.UserName == userName {
		return true
	}
	var countNum int
	model.Db.Model(&model.WorkflowRecord{}).Select("count(*) as count_num").
		Where("workflow_id = ? and assignee_user_name = ?", workflowId, userName).First(&countNum)
	if countNum > 0 {
		return true
	}
	return false
}
