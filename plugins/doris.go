package plugins

import (
	"horizon/model"
	"vitess.io/vitess/go/vt/sqlparser"
)

/**
1. 审核
2. 执行
*/

type DorisPlugin struct {
}

type AuditResult struct {
	Statement   string
	AuditStatus model.WorkflowSqlAuditStatus
	AuditLevel  model.WorkflowSqlAuditLevel
	AuditMsg    string
}

func (dp *DorisPlugin) Audit(sqlContent string) ([]AuditResult, error) {
	var auditResults []AuditResult
	pieces, err := dp.Parse(sqlContent)
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
