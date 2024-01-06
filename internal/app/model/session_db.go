package model

import "github.com/google/uuid"

type Session struct {
	ID       uint      `gorm:"primary_key"`
	UserID   uuid.UUID `gorm:"ForeignKey:ID"`
	Language string    `gorm:"type:varchar(100)"`
	Scene    string    `gorm:"type:varchar(255)"`
	Detail   string    `gorm:"type:varchar(500)"`
	Name     string    `gorm:"type:varchar(100)"`
	//DialogueCount int       `gorm:"type:int"`
}
