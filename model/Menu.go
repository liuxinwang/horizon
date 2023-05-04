package model

import (
	"gorm.io/datatypes"
	"time"
)

type Menu struct {
	ID             uint             `gorm:"primaryKey;comment:主键ID" json:"id"`
	Name           string           `gorm:"type:varchar(20);not null;comment:名称" json:"name"`
	ParentId       uint             `gorm:"comment:上级ID" json:"parentId"`
	Meta           datatypes.JSON   `gorm:"type:json;not null;comment:元信息" json:"meta"`
	Component      string           `gorm:"type:varchar(50);not null;comment:组件" json:"component"`
	Redirect       string           `gorm:"type:varchar(255);comment:重定向" json:"redirect"`
	Path           string           `gorm:"type:varchar(255);comment:路径" json:"path"`
	ActionData     datatypes.JSON   `gorm:"type:json;null;comment:按钮操作数据" json:"actionData"`
	ActionList     datatypes.JSON   `gorm:"type:json;null;comment:按钮操作数据列表" json:"actionList"`
	CreatedAt      time.Time        `gorm:"type:datetime;not null;default:current_timestamp;comment:创建时间" json:"createdAt"`
	UpdatedAt      time.Time        `gorm:"type:datetime;not null;default:current_timestamp on update current_timestamp;comment:修改时间" json:"updatedAt"`
	RolePermission []RolePermission `gorm:"foreignKey:MenuId"`
}
