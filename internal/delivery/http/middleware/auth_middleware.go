package middleware

import (
	"go-rbac-starter/internal/model"
	"go-rbac-starter/internal/usecase"
	"slices"

	"github.com/gofiber/fiber/v2"
)

func NewAuth(userUseCase *usecase.UserUseCase) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		request := &model.VerifyUserRequest{Token: ctx.Get("Authorization", "NOT_FOUND")}
		userUseCase.Log.Debugf("Authorization : %s", request.Token)

		auth, err := userUseCase.Verify(ctx.UserContext(), request)
		if err != nil {
			userUseCase.Log.Warnf("Failed find user by token : %+v", err)
			return fiber.ErrUnauthorized
		}

		userUseCase.Log.Debugf("User : %+v", auth.ID)
		ctx.Locals("auth", auth)
		return ctx.Next()
	}
}

func GetUser(ctx *fiber.Ctx) *model.Auth {
	return ctx.Locals("auth").(*model.Auth)
}

func RequirePermission(permission string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		auth := GetUser(ctx)
		if slices.Contains(auth.Permissions, permission) {
			return ctx.Next()
		}
		return fiber.ErrForbidden
	}
}

func RequireRole(role string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		auth := GetUser(ctx)
		if auth.RoleName == role {
			return ctx.Next()
		}
		return fiber.ErrForbidden
	}
}
