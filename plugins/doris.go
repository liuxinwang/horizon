package plugins

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"horizon/config"
	"horizon/model"
	"horizon/utils"
	"vitess.io/vitess/go/vt/sqlparser"
)

/**
1. 审核
2. 执行
*/

type DorisPlugin struct {
	Instance *model.Instance
}

func (dp *DorisPlugin) Audit(workflow *model.Workflow) ([]AuditResult, error) {
	var auditResults []AuditResult
	pieces, err := dp.Parse(workflow.SqlContent)
	if err != nil {
		return nil, err
	}
	for _, piece := range pieces {
		auditResults = append(auditResults, AuditResult{
			Statement:   piece,
			AuditStatus: model.WorkflowSqlAuditStatusPassed,
			AuditLevel:  model.WorkflowSqlAuditLevelSuccess,
			AuditMsg:    "",
		})
	}
	return auditResults, nil
}

func (dp *DorisPlugin) Parse(sql string) ([]string, error) {
	var piecesNotComment []string
	pieces, err := sqlparser.SplitStatementToPieces(sql)
	if err != nil {
		return nil, err
	}
	for _, piece := range pieces {
		rs, _ := sqlparser.SplitMarginComments(piece)
		piecesNotComment = append(piecesNotComment, rs)
	}
	return piecesNotComment, nil
}

func (dp *DorisPlugin) Execute(workflow *model.Workflow) error {
	// 迭代 WorkflowSqlDetail
	rows, err := model.Db.Model(&model.WorkflowSqlDetail{}).
		Where("workflow_id = ?", workflow.ID).Order("serial_number asc").Rows()
	defer rows.Close()
	if err != nil {
		return err
	}

	dsn := "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=1s"
	dsn = fmt.Sprintf(dsn, dp.Instance.User,
		utils.DecryptAES([]byte(config.Conf.General.SecretKey), dp.Instance.Password),
		dp.Instance.Ip, dp.Instance.Port, workflow.DbName)
	Db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	for rows.Next() {
		var workflowSqlDetail model.WorkflowSqlDetail
		model.Db.ScanRows(rows, &workflowSqlDetail)
		// 执行SQL
		result := Db.Exec(workflowSqlDetail.Statement)
		if result.Error != nil {
			// 更新状态 workflowSqlDetail failed
			model.Db.Model(&workflowSqlDetail).Updates(model.WorkflowSqlDetail{
				ExecutionStatus: model.WorkflowSqlExecutionStatusFailed,
				ExecutionMsg:    err.Error(),
			})
			return result.Error
		}
		// 更新状态 workflowSqlDetail successfully
		model.Db.Model(&workflowSqlDetail).Updates(model.WorkflowSqlDetail{ExecutionStatus: model.WorkflowSqlExecutionStatusSuccessfully})
	}
	return nil
}
