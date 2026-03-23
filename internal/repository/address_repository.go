package repository

import (
	"go-rbac-starter/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AddressRepository struct {
	Repository[entity.Address]
	Log *logrus.Logger
}

func NewAddressRepository(log *logrus.Logger) *AddressRepository {
	return &AddressRepository{
		Log: log,
	}
}

func (r *AddressRepository) FindByIdAndContactId(db *gorm.DB, address *entity.Address, id, contactId string) error {
	return db.Where("id = ? AND contact_id = ?", id, contactId).Take(address).Error
}

func (r *AddressRepository) FindAllByContactId(db *gorm.DB, contactId string) ([]entity.Address, error) {
	var addresses []entity.Address
	err := db.Where("contact_id = ?", contactId).Find(&addresses).Error
	return addresses, err
}
