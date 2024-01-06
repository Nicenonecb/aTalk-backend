package repository

import (
	"aTalkBackEnd/internal/app/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SessionRepository struct {
	DB *gorm.DB
}

//func (r *SessionRepository) GetAllSessions() ([]model.Session, error) {
//	var sessions []model.Session
//	if err := r.DB.Find(&sessions).Error; err != nil {
//		return nil, err
//	}
//	return sessions, nil
//}

func (r *SessionRepository) CreateSession(session *model.Session) error {
	return r.DB.Create(session).Error
}

func (r *SessionRepository) DeleteSession(id uint) error {
	return r.DB.Delete(&model.Session{}, id).Error
}

func (r *SessionRepository) UpdateSession(session *model.Session) error {
	return r.DB.Save(session).Error
}

func (r *SessionRepository) GetSessionByID(id uint) (*model.Session, error) {
	var session model.Session
	if err := r.DB.Where("id = ?", id).First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *SessionRepository) GetAllSessionsByUserID(userID uuid.UUID) ([]model.Session, error) {
	var sessions []model.Session
	result := r.DB.Where("user_id = ?", userID).Find(&sessions)
	if result.Error != nil {
		return nil, result.Error
	}
	return sessions, nil
}
