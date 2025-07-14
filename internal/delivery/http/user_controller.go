package http

import (
	"errors"
	"strconv"
	"strings"

	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/common/response"
	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/model"
	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/usecase"
	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type UserController struct {
	Log         *zerolog.Logger
	userUseCase usecase.UserUseCase
	redis       *redis.Client
	validator   *validator.Validate
}

func NewUserController(userUseCase usecase.UserUseCase, log *zerolog.Logger, redis *redis.Client, validator *validator.Validate) *UserController {
	return &UserController{
		Log:         log,
		userUseCase: userUseCase,
		redis:       redis,
		validator:   validator,
	}
}

func (c *UserController) Register(ctx *fiber.Ctx) error {
	request := new(model.RegisterUserRequest)

	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Error().Msgf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}

	if err := c.validator.Struct(request); err != nil {
		return utils.SendValidatorErrorResponse(ctx, err)
	}

	data, err := c.userUseCase.Create(ctx.Context(), request)
	if err != nil {
		c.Log.Warn().Msgf("Failed to create user: %v", err)
		return err
	}

	return utils.SendSuccessResponse(ctx, response.Success, data)
}

func (c *UserController) Login(ctx *fiber.Ctx) error {
	request := new(model.LoginUserRequest)

	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Error().Msgf("Failed to parse request body: %v", err)

		return fiber.ErrBadRequest
	}

	if err := c.validator.Struct(request); err != nil {
		return utils.SendValidatorErrorResponse(ctx, err)
	}

	data, err := c.userUseCase.Login(ctx.Context(), request)
	if err != nil {
		c.Log.Warn().Msgf("Failed to login user: %v", err)
		return err
	}

	return utils.SendSuccessResponse(ctx, response.Success, data)
}

func (h *UserController) GetCurrentUser(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	user, err := h.userUseCase.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound("")) {
			return utils.SendErrorResponse(c, response.BadRequest, err.Error())
		}
		return utils.SendErrorResponse(c, response.ServerError, err.Error())
	}
	return utils.SendSuccessResponse(c, response.Success, user)
}

func (h *UserController) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.userUseCase.GetAllUsers()
	if err != nil {
		return utils.SendErrorResponse(c, response.ServerError, err.Error())
	}
	return utils.SendSuccessResponse(c, response.Success, users)
}

func (h *UserController) GetUserByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.SendErrorResponse(c, response.BadRequest)
	}

	user, err := h.userUseCase.GetUserByID(uint(id))
	if err != nil {
		if errors.Is(err, utils.ErrNotFound("")) {
			return utils.SendErrorResponse(c, response.BadRequest, err.Error())
		}
		return utils.SendErrorResponse(c, response.ServerError, err.Error())
	}
	return utils.SendSuccessResponse(c, response.Success, "Pengguna berhasil diambil", user)
}

func (h *UserController) UpdateUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.SendErrorResponse(c, response.BadRequest, "ID pengguna tidak valid")
	}

	// Otorisasi: Hanya user yang sedang login atau admin yang bisa mengupdate dirinya/user lain.
	// Contoh sederhana: hanya user itu sendiri atau jika ada mekanisme admin.
	// loggedInUserID := c.Locals("userID").(uint)
	// Jika user mencoba update user lain DAN bukan admin (jika ada peran admin), tolak.
	// Untuk demo ini, kita izinkan user yang login mengupdate dirinya sendiri.
	// Jika Anda ingin hanya admin yang bisa update user lain, tambahkan middleware `IsAdmin`
	// sebelum endpoint ini, atau tambahkan logika role di handler ini.
	// if uint(id) != loggedInUserID /* && !c.Locals("isAdmin").(bool) <-- contoh jika ada middleware admin */ {
	// 	return utils.SendErrorResponse(c, http.StatusForbidden, "Anda tidak memiliki izin untuk memperbarui pengguna ini")
	// }

	var req model.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, response.BadRequest, "Invalid request body")
	}

	if err := h.validator.Struct(req); err != nil {
		return utils.SendValidatorErrorResponse(c, err)
	}

	user, err := h.userUseCase.UpdateUser(uint(id), req.Username, req.Email, req.Password, req.Role)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound("")) {
			return utils.SendErrorResponse(c, response.BadRequest, err.Error())
		}
		if errors.Is(err, utils.ErrValidation("")) {
			return utils.SendErrorResponse(c, response.BadRequest, err.Error())
		}
		if strings.Contains(err.Error(), "digunakan") { // Untuk pesan duplikasi
			return utils.SendErrorResponse(c, response.BadRequest, err.Error())
		}
		return utils.SendErrorResponse(c, response.ServerError, err.Error())
	}
	return utils.SendSuccessResponse(c, response.Success, user)
}

func (h *UserController) DeleteUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.SendErrorResponse(c, response.BadRequest, "ID pengguna tidak valid")
	}

	// Otorisasi: Hanya admin yang boleh menghapus user lain.
	// Atau user hanya boleh menghapus akunnya sendiri.
	// loggedInUserID := c.Locals("userID").(uint)
	// if uint(id) != loggedInUserID  {
	// 	return utils.SendErrorResponse(c, http.StatusForbidden, "Anda tidak memiliki izin untuk menghapus pengguna ini")
	// }

	err = h.userUseCase.DeleteUser(uint(id))
	if err != nil {
		if errors.Is(err, utils.ErrNotFound("")) {
			return utils.SendErrorResponse(c, response.BadRequest, err.Error())
		}
		return utils.SendErrorResponse(c, response.ServerError, err.Error())
	}
	return utils.SendSuccessResponse(c, response.Success)
}
