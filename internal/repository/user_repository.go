package repository

import (
	"go-rbac-starter/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepository struct {
	Repository[entity.User]
	Log *logrus.Logger
}

func NewUserRepository(log *logrus.Logger) *UserRepository {
	return &UserRepository{
		Log: log,
	}
}

func (r *UserRepository) FindByEmail(db *gorm.DB, user *entity.User, email string) error {
	return db.Where("email = ?", email).First(user).Error
}

func (r *UserRepository) FindByToken(db *gorm.DB, user *entity.User, token string) error {
	return db.Where("token = ?", token).First(user).Error
}
