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

func (r *UserRepository) FindByTokenWithRole(db *gorm.DB, user *entity.User, token string) error {
	return db.Preload("Role.Permissions").Where("token = ?", token).First(user).Error
}

func (r *UserRepository) FindByIdWithRole(db *gorm.DB, user *entity.User, id int) error {
	return db.Preload("Role").Where("id = ?", id).First(user).Error
}

func (r *UserRepository) FindAllWithRole(db *gorm.DB, offset, limit int) ([]entity.User, int64, error) {
	var users []entity.User
	var total int64

	if err := db.Model(&entity.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Preload("Role").Offset(offset).Limit(limit).Find(&users).Error
	return users, total, err
}
