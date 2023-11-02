package repository

import (
	"aTalkBackEnd/internal/app/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	if err := r.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(user *model.User) error {
	return r.DB.Create(user).Error
}
