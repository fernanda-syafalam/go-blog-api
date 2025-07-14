package config

import (
	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/delivery/http"
	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/delivery/http/route"
	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/repository"
	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/usecase"
	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/utils"
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
	utils.InitValidator()

	// Register Repository
	userRespository := repository.NewUserRepository(config.Log, config.DB)
	postRepository := repository.NewPostRepository(config.DB)
	categoryRepository := repository.NewCategoryRepository(config.DB)
	commentRepository := repository.NewCommentRepository(config.DB)

	// Register UseCase
	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, config.Validate, userRespository, config.Config)
	postUseCase := usecase.NewPostUseCase(postRepository,categoryRepository, config.Validate)
	categoryUseCase := usecase.NewCategoryUseCase(categoryRepository, config.Validate)
	commentUseCase := usecase.NewCommentUseCase(commentRepository,postRepository, config.Validate)

	// Register Controller
	userController := http.NewUserController(userUseCase, config.Log, config.Redis, config.Validate)
	postController := http.NewPostController(postUseCase, config.Validate)
	categoryController := http.NewCategoryController(categoryUseCase, config.Validate)
	commentController := http.NewCommentController(commentUseCase, config.Validate)

	routeConfig := route.RouteConfig{
		App:            config.App,
		Config:         config.Config,
		UserController: userController,
		PostController: postController,
		CategoryController: categoryController,
		CommentController: commentController,
	}

	routeConfig.Setup()
}

