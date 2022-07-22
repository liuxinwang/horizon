package model

import "time"

type Instance struct {
	ID         uint         `gorm:"primaryKey;comment:主键ID"`
	InstId     string       `gorm:"type:varchar(20);uniqueIndex:uniq_inst_id;not null;comment:实例ID"`
	Name       string       `gorm:"type:varchar(50);not null;comment:实例名称"`
	Role       string       `gorm:"type:enum('Master', 'Slave');default:'Master';not null;comment:实例角色"`
	Ip         string       `gorm:"type:varchar(20);not null;comment:实例IP"`
	Port       uint16       `gorm:"not null;comment:实例端口"`
	User       string       `gorm:"type:varchar(50);not null;comment:用户名"`
	Password   string       `gorm:"type:varchar(100);not null;comment:密码"`
	Version    string       `gorm:"type:varchar(50);not null;comment:实例版本"`
	Status     string       `gorm:"type:enum('Start', 'Stop', 'Running', 'Error');default:'Running';not null;comment:实例状态"`
	InspStatus string       `gorm:"type:enum('Enabled', 'Disabled');not null;default:'Enabled';comment:巡检状态；Disabled：关闭；Enabled：开启"`
	CreatedAt  time.Time    `gorm:"type:datetime;not null;default:current_timestamp;comment:创建时间"`
	UpdatedAt  time.Time    `gorm:"type:datetime;not null;default:current_timestamp on update current_timestamp;comment:修改时间"`
	Inspection []Inspection `gorm:"foreignKey:InstId;references:InstId"`
	InstMetric []InstMetric `gorm:"foreignKey:InstId;references:InstId"`
}
