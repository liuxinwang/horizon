package model

import "time"
import "github.com/shopspring/decimal"

type Score struct {
	ID        uint            `gorm:"primaryKey;comment:主键ID"`
	InspId    string          `gorm:"type:varchar(50);not null;comment:巡检ID"`
	Metric    string          `gorm:"type:varchar(50);not null;comment:指标Key"`
	Value     decimal.Decimal `gorm:"type:decimal(10,2);not null;comment:巡检值"`
	Deduction int8            `gorm:"type:int;not null;comment:扣分"`
	Message   string          `gorm:"type:varchar(255);not null;default:'';comment:信息"`
	CreatedAt time.Time       `gorm:"type:datetime;not null;default:current_timestamp;comment:创建时间"`
	UpdatedAt time.Time       `gorm:"type:datetime;not null;default:current_timestamp on update current_timestamp;comment:修改时间"`
}
