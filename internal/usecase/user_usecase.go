package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/alexedwards/argon2id"
	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/entity"
	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/model"
	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/model/converter"
	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/repository"
	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/knadh/koanf"
	"github.com/rs/zerolog"

	"gorm.io/gorm"
)

type userUseCaseImpl struct {
	DB             *gorm.DB
	Log            *zerolog.Logger
	cfg            *koanf.Koanf
	Validate       *validator.Validate
	UserRepository repository.UserRepository
}

type UserUseCase interface {
	Verify(ctx context.Context, request *model.VerifyUserRequest) (*model.Auth, error)
	Create(ctx context.Context, request *model.RegisterUserRequest) (*model.UserResponse, error)
	Login(ctx context.Context, request *model.LoginUserRequest) (*model.UserResponse, error)
	GetAllUsers() ([]entity.User, error)
	GetUserByID(id uint) (*entity.User, error)
	UpdateUser(id uint, username, email, password, role *string) (*entity.User, error)
	DeleteUser(id uint) error
}

var (
	JwtSecret string
	JwtExpire int
)

func NewUserUseCase(db *gorm.DB, log *zerolog.Logger, validate *validator.Validate, UserRepository repository.UserRepository, config *koanf.Koanf) *userUseCaseImpl {
	JwtExpire = config.Int("jwt.expiration")
	JwtSecret = config.String("jwt.secret")

	return &userUseCaseImpl{
		DB:             db,
		Log:            log,
		cfg:            config,
		Validate:       validate,
		UserRepository: UserRepository,
	}
}

func (c *userUseCaseImpl) Verify(ctx context.Context, request *model.VerifyUserRequest) (*model.Auth, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warn().Msgf("Invalid reques body : %v", err)
		return nil, fiber.ErrBadRequest
	}

	user := new(entity.User)
	if err := c.UserRepository.FindByToken(user, request.Token); err != nil {
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

func (c *userUseCaseImpl) Create(ctx context.Context, request *model.RegisterUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warn().Msgf("Invalid reques body : %v", err)
		return nil, fiber.ErrBadRequest
	}

	total, err := c.UserRepository.CountByEmail(request.Email)
	if err != nil {
		c.Log.Warn().Msgf("Failed to count user by id : %v", err)
		return nil, fiber.ErrInternalServerError
	}

	if total > 0 {
		c.Log.Warn().Msgf("User with id %s already exists", request.Email)
		return nil, fiber.ErrConflict
	}

	password, err := argon2id.CreateHash(request.Password, argon2id.DefaultParams)
	if err != nil {
		c.Log.Warn().Msgf("Failed to hash password : %v", err)
		return nil, fiber.ErrInternalServerError
	}

	user := &entity.User{
		Email:        request.Email,
		PasswordHash: string(password),
		Username:     request.Username,
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

func (c *userUseCaseImpl) Login(ctx context.Context, request *model.LoginUserRequest) (*model.UserResponse, error) {

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warn().Msgf("Invalid request body: %v", err)

		return nil, fiber.ErrBadRequest
	}

	user, err := c.UserRepository.FindByEmail(request.Email)
	if err != nil {
		c.Log.Warn().Msgf("User not found: %v", err)

		return nil, fiber.ErrNotFound
	}
	fmt.Println(user)

	match, err := argon2id.ComparePasswordAndHash(request.Password, user.PasswordHash)
	if err != nil || !match {
		c.Log.Warn().Msg("Password mismatch")

		return nil, fiber.ErrUnauthorized
	}

	accessToken, err := utils.GenerateToken(user.ID, user.Role, JwtSecret, JwtExpire)
	if err != nil {
		c.Log.Error().Msgf("Failed to generate access token: %v", err)

		return nil, fiber.ErrInternalServerError
	}
	return &model.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Token:    accessToken,
	}, nil
}

func (s *userUseCaseImpl) GetAllUsers() ([]entity.User, error) {
	users, err := s.UserRepository.FindAll()
	if err != nil {
		return nil, errors.New("Gagal mengambil semua pengguna: " + err.Error())
	}
	return users, nil
}

func (s *userUseCaseImpl) GetUserByID(id uint) (*entity.User, error) {
	user, err := s.UserRepository.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound("Pengguna")
		}
		return nil, errors.New("Gagal mengambil pengguna: " + err.Error())
	}
	return user, nil
}

func (s *userUseCaseImpl) UpdateUser(id uint, username, email, password, role *string) (*entity.User, error) {
	user, err := s.UserRepository.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound("Pengguna")
		}
		return nil, errors.New("Gagal menemukan pengguna: " + err.Error())
	}

	user.Username = *username

	if email != nil && *email != "" {
		if !strings.EqualFold(user.Email, *email) {
			existingUser, err := s.UserRepository.FindByEmail(*email)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("Gagal memeriksa email: " + err.Error())
			}
			if existingUser != nil && existingUser.ID != 0 && existingUser.ID != user.ID {
				return nil, utils.ErrValidation("Email '" + *email + "' sudah digunakan")
			}
			user.Email = *email
		}
	}

	if password != nil && *password != "" {
		hashedPassword, err := argon2id.CreateHash(*password, argon2id.DefaultParams)
		if err != nil {
			return nil, errors.New("Gagal mengenkripsi password: " + err.Error())
		}
		user.PasswordHash = string(hashedPassword)
	}

	if role != nil && *role != "" {
		switch entity.UserRole(*role) {
		case entity.UserRoleReader, entity.UserRoleAuthor, entity.UserRoleAdmin:
			user.Role = entity.UserRole(*role)
		default:
			return nil, utils.ErrValidation("Role pengguna tidak valid: " + *role)
		}
	}

	err = s.UserRepository.Update(s.DB, user)
	if err != nil {
		return nil, errors.New("Gagal memperbarui pengguna: " + err.Error())
	}
	return user, nil
}

func (s *userUseCaseImpl) DeleteUser(id uint) error {
	_, err := s.UserRepository.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrNotFound("Pengguna")
		}
		return errors.New("Gagal menemukan pengguna: " + err.Error())
	}

	err = s.UserRepository.Delete(id)
	if err != nil {
		return errors.New("Gagal menghapus pengguna: " + err.Error())
	}
	return nil
}
