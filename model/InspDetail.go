package model

import (
	"gorm.io/datatypes"
	"time"
)

type InspDetail struct {
	ID        uint           `gorm:"primaryKey;comment:主键ID"`
	InspId    string         `gorm:"type:varchar(50);not null;comment:巡检ID"`
	Metric    string         `gorm:"type:varchar(50);not null;comment:指标Key"`
	Result    datatypes.JSON `gorm:"type:json;not null;comment:结果"`
	CreatedAt time.Time      `gorm:"type:datetime;not null;default:current_timestamp;comment:创建时间"`
	UpdatedAt time.Time      `gorm:"type:datetime;not null;default:current_timestamp on update current_timestamp;comment:修改时间"`
}
