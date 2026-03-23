package repository

import (
	"go-rbac-starter/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RoleRepository struct {
	Repository[entity.Role]
	Log *logrus.Logger
}

func NewRoleRepository(log *logrus.Logger) *RoleRepository {
	return &RoleRepository{
		Log: log,
	}
}

func (r *RoleRepository) FindByName(db *gorm.DB, role *entity.Role, name string) error {
	return db.Where("name = ?", name).Take(role).Error
}

func (r *RoleRepository) FindByIdWithPermissions(db *gorm.DB, role *entity.Role, id int) error {
	return db.Preload("Permissions").Where("id = ?", id).Take(role).Error
}

func (r *RoleRepository) FindAllWithPermissions(db *gorm.DB) ([]entity.Role, error) {
	var roles []entity.Role
	err := db.Preload("Permissions").Find(&roles).Error
	return roles, err
}
