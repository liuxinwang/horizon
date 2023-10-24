package plugins

import (
	"database/sql"
	"fmt"
	"horizon/config"
	"horizon/model"
	"horizon/utils"
	"strings"
	"vitess.io/vitess/go/vt/sqlparser"
)

/**
1. 审核
2. 执行
*/

type MySQLPlugin struct {
	Instance *model.Instance
}

func (mp *MySQLPlugin) Audit(workflow *model.Workflow) ([]AuditResult, error) {
	var auditResults []AuditResult

	dsn := "%s:%s@tcp(%s:%d)/"
	dsn = fmt.Sprintf(dsn, config.Conf.GoInception.User, config.Conf.GoInception.Password,
		config.Conf.GoInception.Host, config.Conf.GoInception.Port)
	db, err := sql.Open("mysql", dsn)
	defer db.Close()

	connInfo := fmt.Sprintf(
		`/*--user=%s;--password=%s;--host=%s;--port=%d;--check=1;*/`,
		mp.Instance.User,
		utils.DecryptAES([]byte(config.Conf.General.SecretKey), mp.Instance.Password),
		mp.Instance.Ip, mp.Instance.Port)

	auditSql := fmt.Sprintf(`%s
    inception_magic_start;
	use %s;
    %s;
    inception_magic_commit;`, connInfo, workflow.DbName, strings.TrimRight(workflow.SqlContent, ";"))

	rows, err := db.Query(auditSql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var orderId, stage, errorLevel, stageStatus, errorMessage, rowSql, affectedRows, sequence, backupDbname, executeTime, sqlsha1, backupTime []uint8
		err = rows.Scan(&orderId, &stage, &errorLevel, &stageStatus, &errorMessage, &rowSql, &affectedRows, &sequence, &backupDbname, &executeTime, &sqlsha1, &backupTime)
		var auditStatus model.WorkflowSqlAuditStatus
		var auditLevel model.WorkflowSqlAuditLevel
		if string(stageStatus) == "Audit Completed" {
			auditStatus = model.WorkflowSqlAuditStatusPassed
		} else {
			auditStatus = model.WorkflowSqlAuditStatusFailed
		}
		if string(errorLevel) == "0" {
			auditLevel = model.WorkflowSqlAuditLevelSuccess
		} else if string(errorLevel) == "1" {
			auditLevel = model.WorkflowSqlAuditLevelWarning
		} else {
			auditLevel = model.WorkflowSqlAuditLevelError
		}
		auditResults = append(auditResults, AuditResult{
			Statement:   string(rowSql),
			AuditStatus: auditStatus,
			AuditLevel:  auditLevel,
			AuditMsg:    string(errorMessage),
		})
	}
	return auditResults, nil
}

func (mp *MySQLPlugin) Parse(sql string) ([]string, error) {
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

func (mp *MySQLPlugin) Execute(workflow *model.Workflow) error {
	dsn := "%s:%s@tcp(%s:%d)/"
	dsn = fmt.Sprintf(dsn, config.Conf.GoInception.User, config.Conf.GoInception.Password,
		config.Conf.GoInception.Host, config.Conf.GoInception.Port)
	db, err := sql.Open("mysql", dsn)
	defer db.Close()

	connInfo := fmt.Sprintf(
		`/*--user=%s;--password=%s;--host=%s;--port=%d;--execute=1;--ignore-warnings=1;--sleep=200;--sleep_rows=100*/`,
		mp.Instance.User,
		utils.DecryptAES([]byte(config.Conf.General.SecretKey), mp.Instance.Password),
		mp.Instance.Ip, mp.Instance.Port)

	executeSql := fmt.Sprintf(`%s
    inception_magic_start;
	use %s;
    %s;
    inception_magic_commit;`, connInfo, workflow.DbName, strings.TrimRight(workflow.SqlContent, ";"))

	rows, err := db.Query(executeSql)
	if err != nil {
		return err
	}
	defer rows.Close()

	var rsErr error

	for rows.Next() {
		var orderId, stage, errorLevel, stageStatus, errorMessage, rowSql, affectedRows, sequence, backupDbname, executeTime, sqlsha1, backupTime []uint8
		err = rows.Scan(&orderId, &stage, &errorLevel, &stageStatus, &errorMessage, &rowSql, &affectedRows, &sequence, &backupDbname, &executeTime, &sqlsha1, &backupTime)
		var executeStatus model.WorkflowSqlExecutionStatus
		if string(stageStatus) == "Execute Successfully" {
			executeStatus = model.WorkflowSqlExecutionStatusSuccessfully
		} else {
			executeStatus = model.WorkflowSqlExecutionStatusFailed
		}

		// 更新状态 workflowSqlDetail failed
		model.Db.Model(&model.WorkflowSqlDetail{}).
			Where("workflow_id = ? and serial_number = ?", workflow.ID, string(orderId)).
			Updates(model.WorkflowSqlDetail{
				ExecutionStatus: executeStatus,
				ExecutionMsg:    string(errorMessage),
			})

		if (string(errorLevel) == "1" || string(errorLevel) == "2") &&
			executeStatus != model.WorkflowSqlExecutionStatusSuccessfully {
			rsErr = fmt.Errorf("line %v has error/warning: %v", string(orderId), string(errorMessage))
			break
		}
	}

	return rsErr
}
