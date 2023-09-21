package tasks

import (
	"errors"
	"fmt"
	"horizon/model"
	"horizon/service"
	"time"
	"vitess.io/vitess/go/vt/log"
)

// WorkflowTaskRunning 工单定时任务入口
func WorkflowTaskRunning() {
	ticker := time.NewTicker(time.Second * time.Duration(10))
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			workflows, err := queryScheduledWorkflow()
			if err != nil {
				log.Errorf("query scheduled workflow error, err: %v", err.Error())
			}
			if len(workflows) > 0 {
				for _, workflow := range workflows {
					err = workflowExecute(workflow)
					if err != nil {
						log.Errorf("scheduled workflow execute error, err: %v", err.Error())
					}
				}
			}
		}
	}
}

func queryScheduledWorkflow() ([]*model.Workflow, error) {
	// iter scheduled workflow
	rows, err := model.Db.Model(&model.Workflow{}).
		Where("status = ? and scheduled_at between date_add(now(), interval -1 hour) and now()",
			model.WorkflowStatusScheduledExecution).Order("scheduled_at asc").Rows()
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	var workflows []*model.Workflow
	for rows.Next() {
		var workflow model.Workflow
		model.Db.ScanRows(rows, &workflow)
		workflows = append(workflows, &workflow)
	}
	return workflows, nil
}

func workflowExecute(workflow *model.Workflow) error {
	// 获取实例信息
	var instance model.Instance
	result := model.Db.First(&instance, "inst_id = ?", workflow.InstId)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New(fmt.Sprintf("%v 实例不存在", workflow.InstId))
	}
	// 更新状态 WorkflowStatusExecuting
	result = model.Db.Model(&workflow).Where("id = ?", workflow.ID).Updates(&model.Workflow{Status: model.WorkflowStatusExecuting})
	if result.Error != nil {
		return result.Error
	}

	// 迭代 WorkflowSqlDetail
	rows, err := model.Db.Model(&model.WorkflowSqlDetail{}).
		Where("workflow_id = ?", workflow.ID).Order("serial_number asc").Rows()
	defer rows.Close()
	if err != nil {
		return err
	}
	for rows.Next() {
		var workflowSqlDetail model.WorkflowSqlDetail
		model.Db.ScanRows(rows, &workflowSqlDetail)
		// 执行SQL
		err := service.ExecuteSQL(&instance, workflow.DbName, workflowSqlDetail.Statement)
		if err != nil {
			// 更新状态 workflowSqlDetail failed
			model.Db.Model(&workflowSqlDetail).Updates(model.WorkflowSqlDetail{
				ExecutionStatus: model.WorkflowSqlExecutionStatusFailed,
				ExecutionMsg:    err.Error(),
			})

			// 更新状态 WorkflowStatusExecutionFailed
			result = model.Db.Model(&workflow).Where("id = ?", workflow.ID).Updates(&model.Workflow{Status: model.WorkflowStatusExecutionFailed})
			if result.Error != nil {
				return result.Error
			}
			return err
		}
		// 更新状态 workflowSqlDetail successfully
		model.Db.Model(&workflowSqlDetail).Updates(model.WorkflowSqlDetail{ExecutionStatus: model.WorkflowSqlExecutionStatusSuccessfully})
	}
	// 更新状态 WorkflowStatusFinished
	result = model.Db.Model(&workflow).Where("id = ?", workflow.ID).Updates(&model.Workflow{Status: model.WorkflowStatusFinished})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
