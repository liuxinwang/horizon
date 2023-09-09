package model

import (
	"time"
)

type Project struct {
	ID                   uint                `gorm:"primaryKey;comment:主键ID" json:"id"`
	ProjId               string              `gorm:"type:varchar(20);uniqueIndex:uniq_proj_id;not null;comment:项目ID" json:"projId"`
	Name                 string              `gorm:"type:varchar(50);not null;comment:项目名称" json:"name"`
	Describe             string              `gorm:"type:varchar(255);not null;comment:描述" json:"describe"`
	Status               string              `gorm:"type:enum('Enabled', 'Disabled');default:'Enabled';not null;comment:正常，禁用" json:"status"`
	WorkflowTemplateCode uint                `gorm:"not null;comment:编号" json:"workflowTemplateCode"`
	CreatedAt            time.Time           `gorm:"type:datetime;not null;default:current_timestamp;comment:创建时间" json:"createdAt"`
	UpdatedAt            time.Time           `gorm:"type:datetime;not null;default:current_timestamp on update current_timestamp;comment:修改时间" json:"updatedAt"`
	ProjectDatasources   []ProjectDatasource `gorm:"foreignKey:ProjId;references:ProjId" json:"projectDatasources"`
	ProjectUsers         []ProjectUser       `gorm:"foreignKey:ProjId;references:ProjId" json:"projectUsers"`
}

type ProjectRole struct {
	ID                      string                   `gorm:"type:varchar(50);primaryKey;comment:角色ID" json:"id"`
	Name                    string                   `gorm:"type:varchar(20);uniqueIndex:uniq_name;not null;comment:名称" json:"name"`
	Describe                string                   `gorm:"type:varchar(255);not null;comment:描述" json:"describe"`
	CreatedAt               time.Time                `gorm:"type:datetime;not null;default:current_timestamp;comment:创建时间" json:"createdAt"`
	UpdatedAt               time.Time                `gorm:"type:datetime;not null;default:current_timestamp on update current_timestamp;comment:修改时间" json:"updatedAt"`
	ProjectUsers            []ProjectUser            `gorm:"foreignKey:RoleId;references:ID" json:"projectUsers"`
	WorkflowTemplateDetails []WorkflowTemplateDetail `gorm:"foreignKey:ProjectRoleId;references:ID" json:"workflowTemplateDetails"`
}

type ProjectDatasource struct {
	ID        uint      `gorm:"primaryKey;comment:主键ID" json:"id"`
	ProjId    string    `gorm:"uniqueIndex:uniq_projid_instid;not null;comment:项目ID" json:"projId"`
	InstId    string    `gorm:"uniqueIndex:uniq_projid_instid;not null;comment:实例ID" json:"instId"`
	CreatedAt time.Time `gorm:"type:datetime;not null;default:current_timestamp;comment:创建时间" json:"createdAt"`
	UpdatedAt time.Time `gorm:"type:datetime;not null;default:current_timestamp on update current_timestamp;comment:修改时间" json:"updatedAt"`
}

type ProjectUser struct {
	ID        uint      `gorm:"primaryKey;comment:主键ID" json:"id"`
	ProjId    string    `gorm:"uniqueIndex:uniq_projid_roleid_username;not null;comment:项目ID" json:"projId"`
	RoleId    string    `gorm:"uniqueIndex:uniq_projid_roleid_username;not null;comment:角色ID" json:"roleId"`
	UserName  string    `gorm:"uniqueIndex:uniq_projid_roleid_username;not null;comment:用户名" json:"userName"`
	CreatedAt time.Time `gorm:"type:datetime;not null;default:current_timestamp;comment:创建时间" json:"createdAt"`
	UpdatedAt time.Time `gorm:"type:datetime;not null;default:current_timestamp on update current_timestamp;comment:修改时间" json:"updatedAt"`
}
