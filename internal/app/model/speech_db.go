package model

import (
	"time"
)

type SpeechRecord struct {
	ID              string     `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID          string     `gorm:"type:uuid" json:"user_id"`
	AudioURL        string     `json:"audio_url"`
	TranscribedText string     `json:"transcribed_text"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}
