package service

import (
	"aTalkBackEnd/internal/app/model"
	"aTalkBackEnd/internal/app/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Repo *repository.UserRepository
}

func (s *UserService) Register(user *model.User) error {
	user.ID = uuid.New()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return s.Repo.Create(user)
}

func (s *UserService) Authenticate(username, password string) (bool, error) {
	user, err := s.Repo.FindByUsername(username)
	if err != nil {
		return false, err
	}
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil, nil
}

func (s *UserService) FindByUsername(username string) (*model.User, error) {
	return s.Repo.FindByUsername(username)
}
