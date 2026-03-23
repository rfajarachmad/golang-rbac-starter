package config

import (
	"go-rbac-starter/internal/delivery/http"
	"go-rbac-starter/internal/delivery/http/middleware"
	"go-rbac-starter/internal/delivery/http/route"
	"go-rbac-starter/internal/repository"
	"go-rbac-starter/internal/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	App      *fiber.App
	Log      *logrus.Logger
	Validate *validator.Validate
	Config   *viper.Viper
}

func Bootstrap(config *BootstrapConfig) {
	// setup repositories
	userRepository := repository.NewUserRepository(config.Log)
	contactRepository := repository.NewContactRepository(config.Log)
	addressRepository := repository.NewAddressRepository(config.Log)
	roleRepository := repository.NewRoleRepository(config.Log)

	// setup use cases
	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, config.Validate, userRepository, roleRepository)
	contactUseCase := usecase.NewContactUseCase(config.DB, config.Log, config.Validate, contactRepository)
	addressUseCase := usecase.NewAddressUseCase(config.DB, config.Log, config.Validate, contactRepository, addressRepository)
	roleUseCase := usecase.NewRoleUseCase(config.DB, config.Log, roleRepository)

	// setup controller
	userController := http.NewUserController(userUseCase, config.Log)
	contactController := http.NewContactController(contactUseCase, config.Log)
	addressController := http.NewAddressController(addressUseCase, config.Log)
	roleController := http.NewRoleController(roleUseCase, config.Log)

	// setup middleware
	authMiddleware := middleware.NewAuth(userUseCase)

	routeConfig := route.RouteConfig{
		App:               config.App,
		UserController:    userController,
		ContactController: contactController,
		AddressController: addressController,
		RoleController:    roleController,
		AuthMiddleware:    authMiddleware,
	}
	routeConfig.Setup()
}
