package http

import (
	"errors"

	"strconv"

	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/common/response"
	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/usecase"
	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// CategoryController merepresentasikan handler untuk kategori
type CategoryController struct {
	categoryService usecase.CategoryService
	validator       *validator.Validate
}

func NewCategoryController(categoryService usecase.CategoryService, validator *validator.Validate) *CategoryController {
	return &CategoryController{categoryService: categoryService, validator: validator}
}

type CreateCategoryRequest struct {
	Name string `json:"name" validate:"required,min=2,max=50"`
}

func (h *CategoryController) CreateCategory(c *fiber.Ctx) error {
	var req CreateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, response.BadRequest, "Invalid request body")
	}

	if err := h.validator.Struct(req); err != nil {
		return utils.SendValidatorErrorResponse(c, err)
	}

	category, err := h.categoryService.CreateCategory(req.Name)
	if err != nil {
		if errors.Is(err, utils.ErrValidation("")) {
			return utils.SendErrorResponse(c, response.BadRequest, err.Error())
		}
		return utils.SendErrorResponse(c, response.ServerError, err.Error())
	}

	return utils.SendSuccessResponse(c, response.Success, "Kategori berhasil dibuat", category)
}

func (h *CategoryController) GetAllCategories(c *fiber.Ctx) error {
	categories, err := h.categoryService.GetAllCategories()
	if err != nil {
		return utils.SendErrorResponse(c, response.ServerError, err.Error())
	}

	return utils.SendSuccessResponse(c, response.Success, categories)
}

func (h *CategoryController) GetCategoryByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.SendErrorResponse(c, response.BadRequest, "ID kategori tidak valid")
	}

	category, err := h.categoryService.GetCategoryByID(uint(id))
	if err != nil {
		if errors.Is(err, utils.ErrNotFound("")) {
			return utils.SendErrorResponse(c, response.InvalidRequest, err.Error())
		}
		return utils.SendErrorResponse(c, response.ServerError, err.Error())
	}

	return utils.SendSuccessResponse(c, response.Success, category)
}

// UpdateCategoryRequest merepresentasikan payload request untuk update kategori
type UpdateCategoryRequest struct {
	Name string `json:"name" validate:"required,min=2,max=50"`
}

// UpdateCategory memperbarui kategori
func (h *CategoryController) UpdateCategory(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.SendErrorResponse(c, response.BadRequest, "ID kategori tidak valid")
	}

	var req UpdateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, response.BadRequest, "Invalid request body")
	}

	if err := h.validator.Struct(req); err != nil {
		return utils.SendValidatorErrorResponse(c, err)
	}

	category, err := h.categoryService.UpdateCategory(uint(id), req.Name)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound("")) {
			return utils.SendErrorResponse(c, response.InvalidRequest, err.Error())
		}
		if errors.Is(err, utils.ErrValidation("")) {
			return utils.SendErrorResponse(c, response.BadRequest, err.Error())
		}
		return utils.SendErrorResponse(c, response.ServerError, err.Error())
	}

	return utils.SendSuccessResponse(c, response.Success, category)
}

// DeleteCategory menghapus kategori
func (h *CategoryController) DeleteCategory(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.SendErrorResponse(c, response.BadRequest, "ID kategori tidak valid")
	}

	err = h.categoryService.DeleteCategory(uint(id))
	if err != nil {
		if errors.Is(err, utils.ErrNotFound("")) {
			return utils.SendErrorResponse(c, response.InvalidRequest, err.Error())
		}
		return utils.SendErrorResponse(c, response.ServerError, err.Error())
	}

	return utils.SendSuccessResponse(c, response.Success)
}
