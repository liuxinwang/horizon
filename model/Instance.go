package model

import "time"

type Instance struct {
	ID                 uint                `gorm:"primaryKey;comment:主键ID" json:"id"`
	InstId             string              `gorm:"type:varchar(20);uniqueIndex:uniq_inst_id;not null;comment:实例ID" json:"instId"`
	Name               string              `gorm:"type:varchar(50);not null;comment:实例名称" json:"name"`
	Role               string              `gorm:"type:enum('Master', 'Slave');default:'Master';not null;comment:实例角色" json:"role"`
	Ip                 string              `gorm:"type:varchar(20);not null;comment:实例IP" json:"ip"`
	Port               uint16              `gorm:"not null;comment:实例端口" json:"port"`
	User               string              `gorm:"type:varchar(50);not null;comment:用户名" json:"user"`
	Password           string              `gorm:"type:varchar(100);not null;comment:密码" json:"password"`
	Version            string              `gorm:"type:varchar(50);not null;comment:实例版本" json:"version"`
	Status             string              `gorm:"type:enum('Start', 'Stop', 'Running', 'Error');default:'Running';not null;comment:实例状态" json:"status"`
	InspStatus         string              `gorm:"type:enum('Enabled', 'Disabled');not null;default:'Enabled';comment:巡检状态；Disabled：关闭；Enabled：开启" json:"inspStatus"`
	CreatedAt          time.Time           `gorm:"type:datetime;not null;default:current_timestamp;comment:创建时间" json:"createdAt"`
	UpdatedAt          time.Time           `gorm:"type:datetime;not null;default:current_timestamp on update current_timestamp;comment:修改时间" json:"updatedAt"`
	Inspection         []Inspection        `gorm:"foreignKey:InstId;references:InstId" json:"inspection"`
	InstMetric         []InstMetric        `gorm:"foreignKey:InstId;references:InstId" json:"instMetric"`
	ProjectDatasources []ProjectDatasource `gorm:"foreignKey:InstId;references:InstId" json:"projectDatasources"`
}
