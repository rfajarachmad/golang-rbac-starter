package route

import (
	"go-rbac-starter/internal/delivery/http"
	"go-rbac-starter/internal/delivery/http/middleware"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App               *fiber.App
	UserController    *http.UserController
	ContactController *http.ContactController
	AddressController *http.AddressController
	RoleController    *http.RoleController
	AuthMiddleware    fiber.Handler
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	c.App.Post("/api/users", c.UserController.Register)
	c.App.Post("/api/users/_login", c.UserController.Login)
}

func (c *RouteConfig) SetupAuthRoute() {
	c.App.Use(c.AuthMiddleware)

	// User routes
	c.App.Delete("/api/users", middleware.RequirePermission("user:delete"), c.UserController.Logout)
	c.App.Patch("/api/users/_current", middleware.RequirePermission("user:update"), c.UserController.Update)
	c.App.Get("/api/users/_current", middleware.RequirePermission("user:read"), c.UserController.Current)

	// Contact routes
	c.App.Get("/api/contacts", middleware.RequirePermission("contact:read"), c.ContactController.List)
	c.App.Post("/api/contacts", middleware.RequirePermission("contact:create"), c.ContactController.Create)
	c.App.Put("/api/contacts/:contactId", middleware.RequirePermission("contact:update"), c.ContactController.Update)
	c.App.Get("/api/contacts/:contactId", middleware.RequirePermission("contact:read"), c.ContactController.Get)
	c.App.Delete("/api/contacts/:contactId", middleware.RequirePermission("contact:delete"), c.ContactController.Delete)

	// Address routes
	c.App.Get("/api/contacts/:contactId/addresses", middleware.RequirePermission("address:read"), c.AddressController.List)
	c.App.Post("/api/contacts/:contactId/addresses", middleware.RequirePermission("address:create"), c.AddressController.Create)
	c.App.Put("/api/contacts/:contactId/addresses/:addressId", middleware.RequirePermission("address:update"), c.AddressController.Update)
	c.App.Get("/api/contacts/:contactId/addresses/:addressId", middleware.RequirePermission("address:read"), c.AddressController.Get)
	c.App.Delete("/api/contacts/:contactId/addresses/:addressId", middleware.RequirePermission("address:delete"), c.AddressController.Delete)

	// Admin routes
	c.App.Get("/api/admin/users", middleware.RequirePermission("admin:user:list"), c.UserController.ListAll)
	c.App.Get("/api/admin/users/:userId", middleware.RequirePermission("admin:user:read"), c.UserController.GetAny)
	c.App.Patch("/api/admin/users/:userId/role", middleware.RequirePermission("admin:user:update"), c.UserController.AssignRole)
	c.App.Delete("/api/admin/users/:userId", middleware.RequirePermission("admin:user:delete"), c.UserController.DeleteAny)
	c.App.Get("/api/admin/roles", middleware.RequirePermission("admin:role:manage"), c.RoleController.List)
	c.App.Get("/api/admin/roles/:roleId", middleware.RequirePermission("admin:role:manage"), c.RoleController.Get)
}
