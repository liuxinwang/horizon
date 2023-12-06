package model

import "time"

type DataMigrateJob struct {
	ID           uint      `gorm:"primaryKey;comment:主键ID" json:"id"`
	Name         string    `gorm:"type:varchar(50);not null;comment:任务名称" json:"name"`
	SourceInstId string    `gorm:"type:varchar(20);not null;comment:源实例ID" json:"sourceInstId"`
	SourceDb     string    `gorm:"type:varchar(20);not null;comment:源数据库" json:"sourceDb"`
	TargetInstId string    `gorm:"type:varchar(20);not null;comment:目的实例ID" json:"targetInstId"`
	TargetDb     string    `gorm:"type:varchar(20);not null;comment:目的数据库" json:"targetDb"`
	Status       string    `gorm:"type:enum('Start', 'Stop', 'Running', 'Error', 'Finished');default:'Running';not null;comment:任务状态" json:"status"`
	CreatedAt    time.Time `gorm:"type:datetime;not null;default:current_timestamp;comment:创建时间" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"type:datetime;not null;default:current_timestamp on update current_timestamp;comment:修改时间" json:"updatedAt"`
}

type DataMigrateJobDetail struct {
	ID        uint      `gorm:"primaryKey;comment:主键ID" json:"id"`
	TableName string    `gorm:"type:varchar(50);not null;comment:表名称" json:"tableName"`
	Status    string    `gorm:"type:enum('Start', 'Stop', 'Running', 'Error');default:'Running';not null;comment:同步状态" json:"status"`
	CreatedAt time.Time `gorm:"type:datetime;not null;default:current_timestamp;comment:创建时间" json:"createdAt"`
	UpdatedAt time.Time `gorm:"type:datetime;not null;default:current_timestamp on update current_timestamp;comment:修改时间" json:"updatedAt"`
}
