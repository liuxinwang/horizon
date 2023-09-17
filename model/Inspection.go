package model

import "time"

type Inspection struct {
	ID          uint         `gorm:"primaryKey;comment:主键ID" json:"id"`
	InspId      string       `gorm:"type:varchar(50);uniqueIndex:uniq_insp_id;not null;comment:巡检ID" json:"inspId"`
	InstId      string       `gorm:"type:varchar(20);not null;comment:实例ID" json:"instId"`
	InstName    string       `gorm:"type:varchar(50);not null;comment:实例名称" json:"instName"`
	Score       int8         `gorm:"type:int;not null;default:0;comment:评分" json:"score"`
	Level       string       `gorm:"type:enum('HEALTHY', 'SUBOPTIMAL', 'RISKY', 'CRITICAL', '-');not null;default:'-';comment:健康等级" json:"level"`
	CreatedAt   time.Time    `gorm:"type:datetime;not null;default:current_timestamp;comment:创建时间" json:"createdAt"`
	UpdatedAt   time.Time    `gorm:"type:datetime;not null;default:current_timestamp on update current_timestamp;comment:修改时间" json:"updatedAt"`
	InspDetails []InspDetail `gorm:"foreignKey:InspId;references:InspId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"inspDetails"`
	Scores      []Score      `gorm:"foreignKey:InspId;references:InspId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"scores"`
}
