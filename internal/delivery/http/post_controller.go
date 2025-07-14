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
)

type PostController struct {
	postUseCase usecase.PostUseCase
	validator   *validator.Validate
}

func NewPostController(postUseCase usecase.PostUseCase, validator *validator.Validate) *PostController {
	return &PostController{
		postUseCase: postUseCase,
		validator:   validator,
	}
}

func (c *PostController) CreatePost(ctx *fiber.Ctx) error {
	var request model.CreatePostRequest
	if err := ctx.BodyParser(&request); err != nil {
		return utils.SendErrorResponse(ctx, response.BadRequest)
	}

	if err := c.validator.Struct(request); err != nil {
		return utils.SendValidatorErrorResponse(ctx, err)
	}

	authorID := ctx.Locals("userID").(uint)

	post, err := c.postUseCase.CreatePost(request.Title, request.Content, authorID, request.CategoryNames)
	if err != nil {
		if strings.Contains(err.Error(), "Invalid") || strings.Contains(err.Error(), "has already been taken") {
			return utils.SendErrorResponse(ctx, response.BadRequest, err.Error())
		}

		if errors.Is(err, utils.ErrNotFound("")) || strings.Contains(err.Error(), "already exists") {
			return utils.SendErrorResponse(ctx, response.ServerError, err.Error())
		}
		return utils.SendErrorResponse(ctx, response.ServerError, err.Error())
	}

	return utils.SendSuccessResponse(ctx, response.Success, post)
}

func (h *PostController) GetPostByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.SendErrorResponse(c, response.BadRequest)
	}

	post, err := h.postUseCase.GetPostByID(uint(id))
	if err != nil {
		if errors.Is(err, utils.ErrNotFound("")) {
			return utils.SendErrorResponse(c, response.ServerError, err.Error())
		}
		return utils.SendErrorResponse(c, response.ServerError, err.Error())
	}

	return utils.SendSuccessResponse(c, response.Success, post)
}

func (h *PostController) GetPostBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")

	post, err := h.postUseCase.GetPostBySlug(slug)
	if err != nil {
		if strings.Contains(err.Error(), "Slug tidak valid") {
			return utils.SendErrorResponse(c, response.BadRequest, err.Error())
		}
		if errors.Is(err, utils.ErrNotFound("")) {
			return utils.SendErrorResponse(c, response.ServerError, err.Error())
		}
		return utils.SendErrorResponse(c, response.ServerError, err.Error())
	}

	return utils.SendSuccessResponse(c, response.Success, post)
}

func (h *PostController) GetAllPosts(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	posts, err := h.postUseCase.GetAllPosts(page, limit)
	if err != nil {
		return utils.SendErrorResponse(c, response.ServerError, err.Error())
	}

	return utils.SendSuccessResponse(c, response.Success, posts)
}

func (h *PostController) UpdatePost(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.SendErrorResponse(c, response.BadRequest)
	}

	var req model.UpdatePostRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, response.BadRequest)
	}

	if err := h.validator.Struct(req); err != nil {
		return utils.SendValidatorErrorResponse(c, err)
	}

	authorID := c.Locals("userID").(uint)

	post, err := h.postUseCase.UpdatePost(uint(id), req.Title, req.Content, req.PublishedAt, req.CategoryNames, authorID) // Tambahkan CategoryNames
	if err != nil {
		if errors.Is(err, utils.ErrNotFound("")) {
			return utils.SendErrorResponse(c, response.ServerError, err.Error())
		}
		if errors.Is(err, utils.ErrForbidden("")) {
			return utils.SendErrorResponse(c, response.BadRequest, err.Error())
		}
		if strings.Contains(err.Error(), "tidak valid") || strings.Contains(err.Error(), "sudah terpakai") {
			return utils.SendErrorResponse(c, response.BadRequest, err.Error())
		}
		return utils.SendErrorResponse(c, response.ServerError, err.Error())
	}

	return utils.SendSuccessResponse(c, response.Success, post)
}

func (h *PostController) DeletePost(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.SendErrorResponse(c, response.BadRequest, "ID postingan tidak valid")
	}

	authorID := c.Locals("userID").(uint)

	err = h.postUseCase.DeletePost(uint(id), authorID)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound("")) {
			return utils.SendErrorResponse(c, response.ServerError, err.Error())
		}
		if errors.Is(err, utils.ErrForbidden("")) {
			return utils.SendErrorResponse(c, response.BadRequest, err.Error())
		}
		return utils.SendErrorResponse(c, response.ServerError, err.Error())
	}

	return utils.SendSuccessResponse(c, response.Success)
}
