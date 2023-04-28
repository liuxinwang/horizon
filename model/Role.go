package model

import (
	"time"
)

type Role struct {
	ID              string           `gorm:"type:varchar(50);primaryKey;comment:角色ID" json:"id"`
	Name            string           `gorm:"type:varchar(20);uniqueIndex:uniq_name;not null;comment:名称" json:"name"`
	Describe        string           `gorm:"type:varchar(255);not null;comment:描述" json:"describe"`
	CreatedAt       time.Time        `gorm:"type:datetime;not null;default:current_timestamp;comment:创建时间" json:"createdAt"`
	UpdatedAt       time.Time        `gorm:"type:datetime;not null;default:current_timestamp on update current_timestamp;comment:修改时间" json:"updatedAt"`
	UserRoles       []UserRole       `gorm:"foreignKey:RoleId" json:"userRoles"`
	RolePermissions []RolePermission `gorm:"foreignKey:RoleId" json:"rolePermissions"`
}
