package model

import (
	"time"

	"gorm.io/gorm"
)

// 定义公共模型
type BaseModel struct {
	ID        int32     `gorm:"primary_key"`
	CreateAt  time.Time `gorm:"column:create_at"`
	UpdateAt  time.Time `gorm:"column:update_at"`
	DeleteAt  gorm.DeletedAt
	IsDeleted bool
}

type User struct {
	BaseModel
	Mobile   string     `gorm:"index:idx_mobile;unique;type:varchar(11);not null"`
	Password string     `gorm:"type:varchar(100);not null"`
	NickName string     `gorm:"type:varchar(20)"`
	Birthday *time.Time `gorm:"type:datetime"`
	Gender   string     `gorm:"column:gender;default:male;type:varchar(6) comment 'female表示女，male 表示男'"`
	Role     int        `gorm:"column:role;default:1;type:int comment '1表示普通用户，2表示管理员'"`
}
