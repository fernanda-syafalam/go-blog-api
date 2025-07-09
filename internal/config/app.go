package config

import (
	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/delivery/http"
	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/delivery/http/route"
	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/repository"
	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/knadh/koanf"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	App      *fiber.App
	Log      *zerolog.Logger
	Validate *validator.Validate
	Config   *koanf.Koanf
	Redis    *redis.Client
}

func Boostrap(config *BootstrapConfig) {
	userRespository := repository.NewUserRepository(config.Log)

	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, config.Validate, userRespository)

	userController := http.NewUserController(config.Log, userUseCase, config.Redis)

	routeConfig := route.RouteConfig{
		App:            config.App,
		UserController: userController,
	}

	routeConfig.Setup()
}
