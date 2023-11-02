package service

import (
	"aTalkBackEnd/internal/app/model"
	"aTalkBackEnd/internal/app/repository"
)

type SessionService struct {
	Repo *repository.SessionRepository
}

func (s *SessionService) ListAllSessions() ([]model.Session, error) {
	return s.Repo.GetAllSessions()
}

func (s *SessionService) CreateSession(session *model.Session) error {
	return s.Repo.CreateSession(session)
}

func (s *SessionService) DeleteSession(id uint) error {
	return s.Repo.DeleteSession(id)
}

func (s *SessionService) UpdateSession(session *model.Session) error {
	return s.Repo.UpdateSession(session)
}

func (s *SessionService) GetSessionByID(id uint) (*model.Session, error) {
	return s.Repo.GetSessionByID(id)
}
