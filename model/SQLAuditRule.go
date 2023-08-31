package model

import "time"

type RuleTemplate struct {
	ID        uint      `gorm:"primaryKey;comment:主键ID"`
	Name      string    `gorm:"type:varchar(50);not null;comment:名称"`
	Describe  string    `gorm:"type:varchar(255);not null;comment:描述" json:"describe"`
	DbType    string    `gorm:"type:enum('MySQL');default:'MySQL';not null;comment:类型"`
	CreatedAt time.Time `gorm:"type:datetime;not null;default:current_timestamp;comment:创建时间"`
	UpdatedAt time.Time `gorm:"type:datetime;not null;default:current_timestamp on update current_timestamp;comment:修改时间"`
	// ProjectRules []ProjectRule `gorm:"foreignKey:RuleId" json:"projectRules"`
}
