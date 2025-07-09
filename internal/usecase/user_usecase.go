package usecase

import (
	"context"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/entity"
	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/model"
	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/model/converter"
	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"gorm.io/gorm"
)

var jwtSecretKey = []byte("your-very-secret-key")

type UserUseCase struct {
	DB             *gorm.DB
	Log            *zerolog.Logger
	Validate       *validator.Validate
	UserRepository *repository.UserRepository
}

func NewUserUseCase(db *gorm.DB, log *zerolog.Logger, validate *validator.Validate, userRepository *repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		DB:             db,
		Log:            log,
		Validate:       validate,
		UserRepository: userRepository,
	}
}

func (c *UserUseCase) Verify(ctx context.Context, request *model.VerifyUserRequest) (*model.Auth, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warn().Msgf("Invalid reques body : %v", err)
		return nil, fiber.ErrBadRequest
	}

	user := new(entity.User)
	if err := c.UserRepository.FindByToken(tx, user, request.Token); err != nil {
		c.Log.Warn().Msgf("Failed to find user by token : %v", err)
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warn().Msgf("Failed to commit transaction : %v", err)
		return nil, fiber.ErrInternalServerError
	}

	return &model.Auth{
		ID: user.ID,
	}, nil
}

func (c *UserUseCase) Create(ctx context.Context, request *model.RegisterUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warn().Msgf("Invalid reques body : %v", err)
		return nil, fiber.ErrBadRequest
	}

	total, err := c.UserRepository.CountById(tx, request.ID)
	if err != nil {
		c.Log.Warn().Msgf("Failed to count user by id : %v", err)
		return nil, fiber.ErrInternalServerError
	}

	if total > 0 {
		c.Log.Warn().Msgf("User with id %s already exists", request.ID)
		return nil, fiber.ErrConflict
	}

	password, err := argon2id.CreateHash(request.Password, argon2id.DefaultParams)
	if err != nil {
		c.Log.Warn().Msgf("Failed to hash password : %v", err)
		return nil, fiber.ErrInternalServerError
	}

	user := &entity.User{
		ID:       request.ID,
		Password: string(password),
		Name:     request.Name,
	}

	if err := c.UserRepository.Create(tx, user); err != nil {
		c.Log.Warn().Msgf("Failed to create user : %v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warn().Msgf("Failed to commit transaction : %v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.UserToResponse(user), nil
}

func (c *UserUseCase) Login(ctx context.Context, request *model.LoginUserRequest) (*model.UserResponse, error) {
	tracer := otel.Tracer("user-usecase")
	ctx, span := tracer.Start(ctx, "UserUseCase.Login")
	defer span.End()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warn().Msgf("Invalid request body: %v", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Validation failed")
		return nil, fiber.ErrBadRequest
	}

	user := new(entity.User)
	if err := c.UserRepository.FindById(ctx, c.DB, user, request.ID); err != nil {
		c.Log.Warn().Msgf("User not found: %v", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, "User not found")
		return nil, fiber.ErrNotFound
	}

	match, err := argon2id.ComparePasswordAndHash(request.Password, user.Password)
	if err != nil || !match {
		c.Log.Warn().Msg("Password mismatch")
		span.SetStatus(codes.Error, "Invalid credentials")
		if err != nil {
			span.RecordError(err)
		}
		return nil, fiber.ErrUnauthorized
	}

	accessToken, err := c.generateAccessToken(ctx,user.ID)
	if err != nil {
		c.Log.Error().Msgf("Failed to generate access token: %v", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Token generation failed")
		return nil, fiber.ErrInternalServerError
	}

	span.SetAttributes(
		attribute.String("user.id", user.ID),
		attribute.String("user.name", user.Name),
	)
	span.SetStatus(codes.Ok, "Login success")

	return &model.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Token: accessToken,
	}, nil
}

func (c *UserUseCase) generateAccessToken(ctx context.Context, userID string) (string, error) {
	tracer := otel.Tracer("user-usecase-generate-access-token")
	ctx, span := tracer.Start(ctx, "UserUseCase.generateAccessToken")
	defer span.End()
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(15 * time.Minute).Unix(),
		"iat": time.Now().Unix(),
		"iss": "your-app",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecretKey)
}
