package model

import "time"

type Config struct {
	ID          uint      `gorm:"primaryKey;comment:主键ID"`
	Name        string    `gorm:"type:varchar(50);not null;comment:配置名称"`
	Status      string    `gorm:"type:enum('Start', 'Stop', 'Running', 'Error');default:'Running';not null;comment:实例状态"`
	Description string    `gorm:"type:varchar(100);not null;default:'',comment:配置描述"`
	CreatedAt   time.Time `gorm:"type:datetime;not null;default:current_timestamp;comment:创建时间"`
	UpdatedAt   time.Time `gorm:"type:datetime;not null;default:current_timestamp on update current_timestamp;comment:修改时间"`
}
