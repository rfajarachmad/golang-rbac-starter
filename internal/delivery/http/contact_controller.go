package http

import (
	"go-rbac-starter/internal/delivery/http/middleware"
	"go-rbac-starter/internal/model"
	"go-rbac-starter/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type ContactController struct {
	Log     *logrus.Logger
	UseCase *usecase.ContactUseCase
}

func NewContactController(useCase *usecase.ContactUseCase, logger *logrus.Logger) *ContactController {
	return &ContactController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *ContactController) Create(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := new(model.CreateContactRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}

	request.UserId = auth.ID
	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to create contact : %+v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.ContactResponse]{Data: response})
}

func (c *ContactController) List(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := &model.SearchContactRequest{
		UserId: auth.ID,
		Name:   ctx.Query("name"),
		Email:  ctx.Query("email"),
		Phone:  ctx.Query("phone"),
		Page:   ctx.QueryInt("page", 1),
		Size:   ctx.QueryInt("size", 10),
	}

	responses, paging, err := c.UseCase.Search(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to search contacts : %+v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[[]model.ContactResponse]{
		Data:   responses,
		Paging: paging,
	})
}

func (c *ContactController) Get(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := &model.GetContactRequest{
		ID:     ctx.Params("contactId"),
		UserId: auth.ID,
	}

	response, err := c.UseCase.Get(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to get contact : %+v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.ContactResponse]{Data: response})
}

func (c *ContactController) Update(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := new(model.UpdateContactRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}

	request.ID = ctx.Params("contactId")
	request.UserId = auth.ID

	response, err := c.UseCase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to update contact : %+v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.ContactResponse]{Data: response})
}

func (c *ContactController) Delete(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := &model.DeleteContactRequest{
		ID:     ctx.Params("contactId"),
		UserId: auth.ID,
	}

	response, err := c.UseCase.Delete(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to delete contact : %+v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[bool]{Data: response})
}
