package usecase

import (
	"context"
	"go-rbac-starter/internal/entity"
	"go-rbac-starter/internal/model"
	"go-rbac-starter/internal/model/converter"
	"go-rbac-starter/internal/repository"
	"math"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ContactUseCase struct {
	DB                *gorm.DB
	Log               *logrus.Logger
	Validate          *validator.Validate
	ContactRepository *repository.ContactRepository
}

func NewContactUseCase(db *gorm.DB, logger *logrus.Logger, validate *validator.Validate,
	contactRepository *repository.ContactRepository) *ContactUseCase {
	return &ContactUseCase{
		DB:                db,
		Log:               logger,
		Validate:          validate,
		ContactRepository: contactRepository,
	}
}

func (c *ContactUseCase) Create(ctx context.Context, request *model.CreateContactRequest) (*model.ContactResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	contact := &entity.Contact{
		ID:        uuid.New().String(),
		UserId:    request.UserId,
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
		Phone:     request.Phone,
	}

	if err := c.ContactRepository.Create(tx, contact); err != nil {
		c.Log.Warnf("Failed create contact to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.ContactToResponse(contact), nil
}

func (c *ContactUseCase) Update(ctx context.Context, request *model.UpdateContactRequest) (*model.ContactResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	contact := new(entity.Contact)
	if err := c.ContactRepository.FindByIdAndUserId(tx, contact, request.ID, request.UserId); err != nil {
		c.Log.Warnf("Failed find contact by id : %+v", err)
		return nil, fiber.ErrNotFound
	}

	contact.FirstName = request.FirstName
	contact.LastName = request.LastName
	contact.Email = request.Email
	contact.Phone = request.Phone

	if err := c.ContactRepository.Update(tx, contact); err != nil {
		c.Log.Warnf("Failed save contact : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.ContactToResponse(contact), nil
}

func (c *ContactUseCase) Get(ctx context.Context, request *model.GetContactRequest) (*model.ContactResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	contact := new(entity.Contact)
	if err := c.ContactRepository.FindByIdAndUserId(tx, contact, request.ID, request.UserId); err != nil {
		c.Log.Warnf("Failed find contact by id : %+v", err)
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.ContactToResponse(contact), nil
}

func (c *ContactUseCase) Delete(ctx context.Context, request *model.DeleteContactRequest) (bool, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return false, fiber.ErrBadRequest
	}

	contact := new(entity.Contact)
	if err := c.ContactRepository.FindByIdAndUserId(tx, contact, request.ID, request.UserId); err != nil {
		c.Log.Warnf("Failed find contact by id : %+v", err)
		return false, fiber.ErrNotFound
	}

	if err := c.ContactRepository.Delete(tx, contact); err != nil {
		c.Log.Warnf("Failed delete contact from database : %+v", err)
		return false, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return false, fiber.ErrInternalServerError
	}

	return true, nil
}

func (c *ContactUseCase) Search(ctx context.Context, request *model.SearchContactRequest) ([]model.ContactResponse, *model.PageMetadata, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, nil, fiber.ErrBadRequest
	}

	contacts, total, err := c.ContactRepository.Search(tx, request.UserId, request.Name, request.Email, request.Phone, request.Page, request.Size)
	if err != nil {
		c.Log.Warnf("Failed search contacts from database : %+v", err)
		return nil, nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, nil, fiber.ErrInternalServerError
	}

	responses := make([]model.ContactResponse, 0, len(contacts))
	for i := range contacts {
		responses = append(responses, *converter.ContactToResponse(&contacts[i]))
	}

	paging := &model.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return responses, paging, nil
}
