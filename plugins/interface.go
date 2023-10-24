package plugins

import "horizon/model"

type Plugin interface {
	Audit(workflow *model.Workflow) ([]AuditResult, error)
	Execute(workflow *model.Workflow) error
}

type AuditResult struct {
	Statement   string
	AuditStatus model.WorkflowSqlAuditStatus
	AuditLevel  model.WorkflowSqlAuditLevel
	AuditMsg    string
}
