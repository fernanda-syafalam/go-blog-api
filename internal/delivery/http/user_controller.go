package http

import (
	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/model"
	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel"
)

type UserController struct {
	Log     *zerolog.Logger
	UseCase *usecase.UserUseCase
	redis   *redis.Client
}

func NewUserController(log *zerolog.Logger, useCase *usecase.UserUseCase, redis *redis.Client) *UserController {
	return &UserController{
		Log:     log,
		UseCase: useCase,
		redis:   redis,
	}
}

func (c *UserController) Register(ctx *fiber.Ctx) error {
	request := new(model.RegisterUserRequest)

	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Error().Msgf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.Create(ctx.Context(), request)
	if err != nil {
		c.Log.Warn().Msgf("Failed to create user: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}

func (c *UserController) Login(ctx *fiber.Ctx) error {
	tracer := otel.Tracer("user-controller")
	newCtx, span := tracer.Start(ctx.Context(), "UserController.Login")
	defer span.End()

	request := new(model.LoginUserRequest)

	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Error().Msgf("Failed to parse request body: %v", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Parsing failed")
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.Login(newCtx, request)
	if err != nil {
		c.Log.Warn().Msgf("Failed to login user: %v", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Login failed")
		return err
	}

	span.SetStatus(codes.Ok, "Login success")
	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}
