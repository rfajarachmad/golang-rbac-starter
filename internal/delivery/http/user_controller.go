package http

import (
	"go-rbac-starter/internal/delivery/http/middleware"
	"go-rbac-starter/internal/model"
	"go-rbac-starter/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	Log     *logrus.Logger
	UseCase *usecase.UserUseCase
}

func NewUserController(useCase *usecase.UserUseCase, logger *logrus.Logger) *UserController {
	return &UserController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *UserController) Register(ctx *fiber.Ctx) error {
	request := new(model.RegisterUserRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.Register(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to register user : %+v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}

func (c *UserController) Login(ctx *fiber.Ctx) error {
	request := new(model.LoginUserRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.Login(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to login user : %+v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}

func (c *UserController) Current(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := &model.GetUserRequest{ID: auth.ID}

	response, err := c.UseCase.Current(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to get current user")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}

func (c *UserController) Logout(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := &model.LogoutUserRequest{ID: auth.ID}

	response, err := c.UseCase.Logout(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to logout user")
		return err
	}

	return ctx.JSON(model.WebResponse[bool]{Data: response})
}

func (c *UserController) Update(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := new(model.UpdateUserRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}

	request.ID = auth.ID
	response, err := c.UseCase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to update user")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}

// Admin operations

func (c *UserController) ListAll(ctx *fiber.Ctx) error {
	request := &model.ListUsersRequest{
		Page: ctx.QueryInt("page", 1),
		Size: ctx.QueryInt("size", 10),
	}

	responses, paging, err := c.UseCase.ListAll(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to list users")
		return err
	}

	return ctx.JSON(model.WebResponse[[]model.UserResponse]{Data: responses, Paging: paging})
}

func (c *UserController) GetAny(ctx *fiber.Ctx) error {
	userId, err := ctx.ParamsInt("userId")
	if err != nil {
		c.Log.Warnf("Invalid user id : %+v", err)
		return fiber.ErrBadRequest
	}

	request := &model.GetAnyUserRequest{ID: userId}
	response, err := c.UseCase.GetAny(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to get user")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}

func (c *UserController) AssignRole(ctx *fiber.Ctx) error {
	userId, err := ctx.ParamsInt("userId")
	if err != nil {
		c.Log.Warnf("Invalid user id : %+v", err)
		return fiber.ErrBadRequest
	}

	request := new(model.AssignRoleRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}

	request.UserID = userId
	response, err := c.UseCase.AssignRole(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to assign role")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}

func (c *UserController) DeleteAny(ctx *fiber.Ctx) error {
	userId, err := ctx.ParamsInt("userId")
	if err != nil {
		c.Log.Warnf("Invalid user id : %+v", err)
		return fiber.ErrBadRequest
	}

	request := &model.DeleteAnyUserRequest{ID: userId}
	response, err := c.UseCase.DeleteAny(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to delete user")
		return err
	}

	return ctx.JSON(model.WebResponse[bool]{Data: response})
}
