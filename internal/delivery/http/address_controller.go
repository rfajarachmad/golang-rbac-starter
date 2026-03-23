package http

import (
	"go-rbac-starter/internal/delivery/http/middleware"
	"go-rbac-starter/internal/model"
	"go-rbac-starter/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type AddressController struct {
	Log     *logrus.Logger
	UseCase *usecase.AddressUseCase
}

func NewAddressController(useCase *usecase.AddressUseCase, logger *logrus.Logger) *AddressController {
	return &AddressController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *AddressController) Create(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := new(model.CreateAddressRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}

	request.UserId = auth.ID
	request.ContactId = ctx.Params("contactId")

	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to create address : %+v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.AddressResponse]{Data: response})
}

func (c *AddressController) List(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := &model.ListAddressRequest{
		UserId:    auth.ID,
		ContactId: ctx.Params("contactId"),
	}

	responses, err := c.UseCase.List(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to list addresses : %+v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[[]model.AddressResponse]{Data: responses})
}

func (c *AddressController) Get(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := &model.GetAddressRequest{
		ID:        ctx.Params("addressId"),
		UserId:    auth.ID,
		ContactId: ctx.Params("contactId"),
	}

	response, err := c.UseCase.Get(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to get address : %+v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.AddressResponse]{Data: response})
}

func (c *AddressController) Update(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := new(model.UpdateAddressRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}

	request.ID = ctx.Params("addressId")
	request.UserId = auth.ID
	request.ContactId = ctx.Params("contactId")

	response, err := c.UseCase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to update address : %+v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.AddressResponse]{Data: response})
}

func (c *AddressController) Delete(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := &model.DeleteAddressRequest{
		ID:        ctx.Params("addressId"),
		UserId:    auth.ID,
		ContactId: ctx.Params("contactId"),
	}

	response, err := c.UseCase.Delete(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to delete address : %+v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[bool]{Data: response})
}
