package model

import (
	"gorm.io/datatypes"
	"time"
)

type RolePermission struct {
	ID         uint           `gorm:"primaryKey;comment:主键ID" json:"id"`
	RoleId     string         `gorm:"uniqueIndex:uniq_roleid_menuid;not null;comment:角色ID" json:"roleId"`
	MenuId     uint           `gorm:"uniqueIndex:uniq_roleid_menuid;not null;comment:菜单ID" json:"menuId"`
	ActionData datatypes.JSON `gorm:"type:json;null;comment:按钮可操作数据" json:"actionData"`
	ActionList datatypes.JSON `gorm:"type:json;null;comment:按钮可操作列表数据" json:"actionList"`
	CreatedAt  time.Time      `gorm:"type:datetime;not null;default:current_timestamp;comment:创建时间" json:"createdAt"`
	UpdatedAt  time.Time      `gorm:"type:datetime;not null;default:current_timestamp on update current_timestamp;comment:修改时间" json:"updatedAt"`
	Menu       Menu           `json:"menu"`
	// Permissions datatypes.JSON `gorm:"type:json;null;comment:角色权限"`
}
