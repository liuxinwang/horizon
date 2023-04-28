package model

import "time"

type UserRole struct {
	ID        uint      `gorm:"primaryKey;comment:主键ID"`
	UserId    uint      `gorm:"uniqueIndex:uniq_userid_roleid;not null;comment:用户ID"`
	RoleId    string    `gorm:"uniqueIndex:uniq_userid_roleid;not null;comment:角色ID"`
	CreatedAt time.Time `gorm:"type:datetime;not null;default:current_timestamp;comment:创建时间"`
	UpdatedAt time.Time `gorm:"type:datetime;not null;default:current_timestamp on update current_timestamp;comment:修改时间"`
}
