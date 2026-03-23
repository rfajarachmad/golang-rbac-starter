package test

import (
	"go-rbac-starter/internal/entity"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func ClearAll() {
	ClearAddresses()
	ClearContacts()
	ClearUsers()
}

func ClearUsers() {
	err := db.Where("id is not null").Delete(&entity.User{}).Error
	if err != nil {
		log.Fatalf("Failed clear user data : %+v", err)
	}
}

func ClearContacts() {
	err := db.Where("id is not null").Delete(&entity.Contact{}).Error
	if err != nil {
		log.Fatalf("Failed clear contact data : %+v", err)
	}
}

func ClearAddresses() {
	err := db.Where("id is not null").Delete(&entity.Address{}).Error
	if err != nil {
		log.Fatalf("Failed clear address data : %+v", err)
	}
}

func getRoleId(roleName string) int {
	role := new(entity.Role)
	if err := db.Where("name = ?", roleName).First(role).Error; err != nil {
		log.Fatalf("Failed find role %s : %+v", roleName, err)
	}
	return role.ID
}

func SeedUser(t *testing.T, name, email, password string) *entity.User {
	return SeedUserWithRole(t, name, email, password, "user")
}

func SeedUserWithRole(t *testing.T, name, email, password, roleName string) *entity.User {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.Nil(t, err)

	user := &entity.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Token:    uuid.New().String(),
		RoleId:   getRoleId(roleName),
	}
	err = db.Create(user).Error
	assert.Nil(t, err)
	return user
}

func SeedContact(t *testing.T, userId int) *entity.Contact {
	contact := &entity.Contact{
		ID:        uuid.New().String(),
		UserId:    userId,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Phone:     "08123456789",
	}
	err := db.Create(contact).Error
	assert.Nil(t, err)
	return contact
}

func SeedAddress(t *testing.T, contactId string) *entity.Address {
	address := &entity.Address{
		ID:         uuid.New().String(),
		ContactId:  contactId,
		Street:     "Jalan Test",
		City:       "Jakarta",
		Province:   "DKI Jakarta",
		PostalCode: "12345",
		Country:    "Indonesia",
	}
	err := db.Create(address).Error
	assert.Nil(t, err)
	return address
}

func GetFirstUser(t *testing.T) *entity.User {
	user := new(entity.User)
	err := db.First(user).Error
	assert.Nil(t, err)
	return user
}

func GetFirstContact(t *testing.T, userId int) *entity.Contact {
	contact := new(entity.Contact)
	err := db.Where("user_id = ?", userId).First(contact).Error
	assert.Nil(t, err)
	return contact
}

func GetFirstAddress(t *testing.T, contactId string) *entity.Address {
	address := new(entity.Address)
	err := db.Where("contact_id = ?", contactId).First(address).Error
	assert.Nil(t, err)
	return address
}
