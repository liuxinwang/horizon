package model

import (
	"time"
)

type Task struct {
	ID        uint      `gorm:"primaryKey;comment:主键ID"`
	Name      string    `gorm:"type:varchar(255);uniqueIndex:uniq_name;not null;comment:名称"`
	Args      string    `gorm:"type:varchar(255);not null;comment:参数"`
	CreatedAt time.Time `gorm:"type:datetime;not null;default:current_timestamp;comment:创建时间"`
	UpdatedAt time.Time `gorm:"type:datetime;not null;default:current_timestamp on update current_timestamp;comment:修改时间"`
}
