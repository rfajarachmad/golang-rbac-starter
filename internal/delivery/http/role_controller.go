package http

import (
	"go-rbac-starter/internal/model"
	"go-rbac-starter/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type RoleController struct {
	Log     *logrus.Logger
	UseCase *usecase.RoleUseCase
}

func NewRoleController(useCase *usecase.RoleUseCase, logger *logrus.Logger) *RoleController {
	return &RoleController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *RoleController) List(ctx *fiber.Ctx) error {
	responses, err := c.UseCase.List(ctx.UserContext())
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to list roles")
		return err
	}

	return ctx.JSON(model.WebResponse[[]model.RoleResponse]{Data: responses})
}

func (c *RoleController) Get(ctx *fiber.Ctx) error {
	roleId, err := ctx.ParamsInt("roleId")
	if err != nil {
		c.Log.Warnf("Invalid role id : %+v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.Get(ctx.UserContext(), roleId)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to get role")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.RoleResponse]{Data: response})
}
