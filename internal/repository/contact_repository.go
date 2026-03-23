package repository

import (
	"go-rbac-starter/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ContactRepository struct {
	Repository[entity.Contact]
	Log *logrus.Logger
}

func NewContactRepository(log *logrus.Logger) *ContactRepository {
	return &ContactRepository{
		Log: log,
	}
}

func (r *ContactRepository) FindByIdAndUserId(db *gorm.DB, contact *entity.Contact, id string, userId int) error {
	return db.Where("id = ? AND user_id = ?", id, userId).Take(contact).Error
}

func (r *ContactRepository) Search(db *gorm.DB, userId int, name, email, phone string, page, size int) ([]entity.Contact, int64, error) {
	var contacts []entity.Contact
	var total int64

	query := r.FilterContact(db, userId, name, email, phone)

	if err := query.Model(&entity.Contact{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset((page - 1) * size).Limit(size).Find(&contacts).Error; err != nil {
		return nil, 0, err
	}

	return contacts, total, nil
}

func (r *ContactRepository) FilterContact(db *gorm.DB, userId int, name, email, phone string) *gorm.DB {
	query := db.Where("user_id = ?", userId)

	if name != "" {
		name = "%" + name + "%"
		query = query.Where("first_name ILIKE ? OR last_name ILIKE ?", name, name)
	}

	if email != "" {
		query = query.Where("email ILIKE ?", "%"+email+"%")
	}

	if phone != "" {
		query = query.Where("phone ILIKE ?", "%"+phone+"%")
	}

	return query
}
