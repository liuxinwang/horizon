package model

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey;comment:主键ID"`
	UserName  string    `gorm:"type:varchar(50);not null;comment:用户名""`
	Password  string    `gorm:"type:varchar(100);not null;comment:密码""`
	Status    int32     `gorm:"type:int;not null;default:1;comment:状态；0：无效，1：有效""`
	CreatedAt time.Time `gorm:"type:datetime;not null;default:current_timestamp;comment:创建时间"`
	UpdatedAt time.Time `gorm:"type:datetime;not null;default:current_timestamp on update current_timestamp;comment:修改时间"`
}
