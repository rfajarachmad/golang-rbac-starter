package usecase

import (
	"context"
	"go-rbac-starter/internal/entity"
	"go-rbac-starter/internal/model"
	"go-rbac-starter/internal/model/converter"
	"go-rbac-starter/internal/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RoleUseCase struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	RoleRepository *repository.RoleRepository
}

func NewRoleUseCase(db *gorm.DB, logger *logrus.Logger,
	roleRepository *repository.RoleRepository) *RoleUseCase {
	return &RoleUseCase{
		DB:             db,
		Log:            logger,
		RoleRepository: roleRepository,
	}
}

func (c *RoleUseCase) List(ctx context.Context) ([]model.RoleResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	roles, err := c.RoleRepository.FindAllWithPermissions(tx)
	if err != nil {
		c.Log.Warnf("Failed to list roles : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	responses := make([]model.RoleResponse, len(roles))
	for i := range roles {
		responses[i] = *converter.RoleToResponse(&roles[i])
	}

	return responses, nil
}

func (c *RoleUseCase) Get(ctx context.Context, id int) (*model.RoleResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	role := new(entity.Role)
	if err := c.RoleRepository.FindByIdWithPermissions(tx, role, id); err != nil {
		c.Log.Warnf("Failed find role by id : %+v", err)
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.RoleToResponse(role), nil
}
