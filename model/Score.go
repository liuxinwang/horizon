package model

import "time"
import "github.com/shopspring/decimal"

type Score struct {
	ID        uint            `gorm:"primaryKey;comment:主键ID" json:"id"`
	InspId    string          `gorm:"type:varchar(50);not null;comment:巡检ID" json:"inspId"`
	Metric    string          `gorm:"type:varchar(50);not null;comment:指标Key" json:"metric"`
	Value     decimal.Decimal `gorm:"type:decimal(10,2);not null;comment:巡检值" json:"value"`
	Deduction int8            `gorm:"type:int;not null;comment:扣分" json:"deduction"`
	Message   string          `gorm:"type:varchar(255);not null;default:'';comment:信息" json:"message"`
	CreatedAt time.Time       `gorm:"type:datetime;not null;default:current_timestamp;comment:创建时间" json:"createdAt"`
	UpdatedAt time.Time       `gorm:"type:datetime;not null;default:current_timestamp on update current_timestamp;comment:修改时间" json:"updatedAt"`
}
