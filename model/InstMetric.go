package model

import "time"

type InstMetric struct {
	ID        uint      `gorm:"primaryKey;comment:主键ID"`
	InstId    string    `gorm:"type:varchar(20);uniqueIndex:uniq_inst_id_key;not null;comment:实例ID"`
	Metric    string    `gorm:"type:varchar(50);uniqueIndex:uniq_inst_id_key;not null;comment:指标Key"`
	Status    string    `gorm:"type:enum('Enabled', 'Disabled');not null;default:'Enabled';comment:状态；Disabled：关闭；Enabled：启用"`
	CreatedAt time.Time `gorm:"type:datetime;not null;default:current_timestamp;comment:创建时间"`
	UpdatedAt time.Time `gorm:"type:datetime;not null;default:current_timestamp on update current_timestamp;comment:修改时间"`
}
