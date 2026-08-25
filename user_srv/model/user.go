package model

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        int32     `gorm:"primary_key"`
	CreateAt  time.Time `gorm:"column:create_at"`
	UpdateAt  time.Time `gorm:"column:update_at"`
	DeleteAt  gorm.DeletedAt
	IsDeleted bool
}
