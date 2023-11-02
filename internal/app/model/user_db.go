package model

import (
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key"`
	Username string    `gorm:"type:varchar(100);unique_index"`
	Password string    `gorm:"type:varchar(100)"`
	// ... 其他字段
}
