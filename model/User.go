package model

import "time"

type User struct {
	ID        uint       `gorm:"primaryKey;comment:主键ID" json:"id"`
	UserName  string     `gorm:"type:varchar(50);uniqueIndex:uniq_username;not null;comment:用户名" json:"userName"`
	NickName  string     `gorm:"type:varchar(50);not null;comment:昵称" json:"nickName"`
	Password  string     `gorm:"type:varchar(100);not null;comment:密码" json:"password"`
	Status    string     `gorm:"type:enum('Enabled', 'Disabled');not null;default:'Enabled';comment:状态" json:"status"`
	CreatedAt time.Time  `gorm:"type:datetime;not null;default:current_timestamp;comment:创建时间" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"type:datetime;not null;default:current_timestamp on update current_timestamp;comment:修改时间" json:"updatedAt"`
	UserRoles []UserRole `gorm:"foreignKey:UserId;references:ID" json:"userRoles"`
	// Role      string     `gorm:"type:varchar(50);not null;comment:角色"`
}
