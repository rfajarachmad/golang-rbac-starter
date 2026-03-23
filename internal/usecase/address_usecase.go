package usecase

import (
	"context"
	"go-rbac-starter/internal/entity"
	"go-rbac-starter/internal/model"
	"go-rbac-starter/internal/model/converter"
	"go-rbac-starter/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AddressUseCase struct {
	DB                *gorm.DB
	Log               *logrus.Logger
	Validate          *validator.Validate
	ContactRepository *repository.ContactRepository
	AddressRepository *repository.AddressRepository
}

func NewAddressUseCase(db *gorm.DB, logger *logrus.Logger, validate *validator.Validate,
	contactRepository *repository.ContactRepository, addressRepository *repository.AddressRepository) *AddressUseCase {
	return &AddressUseCase{
		DB:                db,
		Log:               logger,
		Validate:          validate,
		ContactRepository: contactRepository,
		AddressRepository: addressRepository,
	}
}

func (c *AddressUseCase) Create(ctx context.Context, request *model.CreateAddressRequest) (*model.AddressResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	contact := new(entity.Contact)
	if err := c.ContactRepository.FindByIdAndUserId(tx, contact, request.ContactId, request.UserId); err != nil {
		c.Log.Warnf("Failed find contact by id : %+v", err)
		return nil, fiber.ErrNotFound
	}

	address := &entity.Address{
		ID:         uuid.New().String(),
		ContactId:  request.ContactId,
		Street:     request.Street,
		City:       request.City,
		Province:   request.Province,
		PostalCode: request.PostalCode,
		Country:    request.Country,
	}

	if err := c.AddressRepository.Create(tx, address); err != nil {
		c.Log.Warnf("Failed create address to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.AddressToResponse(address), nil
}

func (c *AddressUseCase) Update(ctx context.Context, request *model.UpdateAddressRequest) (*model.AddressResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	contact := new(entity.Contact)
	if err := c.ContactRepository.FindByIdAndUserId(tx, contact, request.ContactId, request.UserId); err != nil {
		c.Log.Warnf("Failed find contact by id : %+v", err)
		return nil, fiber.ErrNotFound
	}

	address := new(entity.Address)
	if err := c.AddressRepository.FindByIdAndContactId(tx, address, request.ID, request.ContactId); err != nil {
		c.Log.Warnf("Failed find address by id : %+v", err)
		return nil, fiber.ErrNotFound
	}

	address.Street = request.Street
	address.City = request.City
	address.Province = request.Province
	address.PostalCode = request.PostalCode
	address.Country = request.Country

	if err := c.AddressRepository.Update(tx, address); err != nil {
		c.Log.Warnf("Failed save address : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.AddressToResponse(address), nil
}

func (c *AddressUseCase) Get(ctx context.Context, request *model.GetAddressRequest) (*model.AddressResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	contact := new(entity.Contact)
	if err := c.ContactRepository.FindByIdAndUserId(tx, contact, request.ContactId, request.UserId); err != nil {
		c.Log.Warnf("Failed find contact by id : %+v", err)
		return nil, fiber.ErrNotFound
	}

	address := new(entity.Address)
	if err := c.AddressRepository.FindByIdAndContactId(tx, address, request.ID, request.ContactId); err != nil {
		c.Log.Warnf("Failed find address by id : %+v", err)
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.AddressToResponse(address), nil
}

func (c *AddressUseCase) Delete(ctx context.Context, request *model.DeleteAddressRequest) (bool, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return false, fiber.ErrBadRequest
	}

	contact := new(entity.Contact)
	if err := c.ContactRepository.FindByIdAndUserId(tx, contact, request.ContactId, request.UserId); err != nil {
		c.Log.Warnf("Failed find contact by id : %+v", err)
		return false, fiber.ErrNotFound
	}

	address := new(entity.Address)
	if err := c.AddressRepository.FindByIdAndContactId(tx, address, request.ID, request.ContactId); err != nil {
		c.Log.Warnf("Failed find address by id : %+v", err)
		return false, fiber.ErrNotFound
	}

	if err := c.AddressRepository.Delete(tx, address); err != nil {
		c.Log.Warnf("Failed delete address from database : %+v", err)
		return false, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return false, fiber.ErrInternalServerError
	}

	return true, nil
}

func (c *AddressUseCase) List(ctx context.Context, request *model.ListAddressRequest) ([]model.AddressResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	contact := new(entity.Contact)
	if err := c.ContactRepository.FindByIdAndUserId(tx, contact, request.ContactId, request.UserId); err != nil {
		c.Log.Warnf("Failed find contact by id : %+v", err)
		return nil, fiber.ErrNotFound
	}

	addresses, err := c.AddressRepository.FindAllByContactId(tx, request.ContactId)
	if err != nil {
		c.Log.Warnf("Failed list addresses from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	responses := make([]model.AddressResponse, 0, len(addresses))
	for i := range addresses {
		responses = append(responses, *converter.AddressToResponse(&addresses[i]))
	}

	return responses, nil
}
