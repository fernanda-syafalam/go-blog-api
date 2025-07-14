package http

import (
	"errors"
	"strconv"
	"strings"

	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/common/response"
	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/usecase"
	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type CommentController struct {
	newCommentUseCase usecase.CommentUseCase
	validator         *validator.Validate
}

func NewCommentController(newCommentUseCase usecase.CommentUseCase, validator *validator.Validate) *CommentController {
	return &CommentController{newCommentUseCase: newCommentUseCase, validator: validator}
}

type CreateCommentRequest struct {
	Content string `json:"content" validate:"required,min=1,max=500"`
}

func (h *CommentController) CreateComment(c *fiber.Ctx) error {
	postID, err := strconv.ParseUint(c.Params("postID"), 10, 32)
	if err != nil {
		return utils.SendErrorResponse(c, response.BadRequest)
	}

	var req CreateCommentRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, response.BadRequest)
	}

	if err := h.validator.Struct(req); err != nil {
		return utils.SendValidatorErrorResponse(c, err)
	}

	authorID := c.Locals("userID").(uint)

	comment, err := h.newCommentUseCase.CreateComment(req.Content, uint(postID), authorID)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound("")) {
			return utils.SendErrorResponse(c, response.BadRequest, err.Error())
		}
		if strings.Contains(err.Error(), "tidak valid") { // Untuk validasi dari service
			return utils.SendErrorResponse(c, response.BadRequest, err.Error())
		}
		return utils.SendErrorResponse(c, response.ServerError, err.Error())
	}

	return utils.SendSuccessResponse(c, response.Success, comment)
}

func (h *CommentController) GetCommentsByPostID(c *fiber.Ctx) error {
	postID, err := strconv.ParseUint(c.Params("postID"), 10, 32)
	if err != nil {
		return utils.SendErrorResponse(c, response.BadRequest)
	}

	comments, err := h.newCommentUseCase.GetCommentsByPostID(uint(postID))
	if err != nil {
		if errors.Is(err, utils.ErrNotFound("")) {
			return utils.SendErrorResponse(c, response.BadRequest, err.Error())
		}
		return utils.SendErrorResponse(c, response.ServerError, err.Error())
	}

	return utils.SendSuccessResponse(c, response.Success, comments)
}

type UpdateCommentRequest struct {
	Content string `json:"content" validate:"required,min=1,max=500"`
}

func (h *CommentController) UpdateComment(c *fiber.Ctx) error {
	commentID, err := strconv.ParseUint(c.Params("commentID"), 10, 32)
	if err != nil {
		return utils.SendErrorResponse(c, response.BadRequest)
	}

	var req UpdateCommentRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, response.BadRequest)
	}

	if err := h.validator.Struct(req); err != nil {
		return utils.SendValidatorErrorResponse(c, err)
	}

	authorID := c.Locals("userID").(uint)

	comment, err := h.newCommentUseCase.UpdateComment(uint(commentID), authorID, req.Content)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound("")) {
			return utils.SendErrorResponse(c, response.BadRequest, err.Error())
		}
		if errors.Is(err, utils.ErrForbidden("")) {
			return utils.SendErrorResponse(c, response.BadRequest, err.Error())
		}
		if strings.Contains(err.Error(), "tidak valid") { // Untuk validasi dari service
			return utils.SendErrorResponse(c, response.BadRequest, err.Error())
		}
		return utils.SendErrorResponse(c, response.ServerError, err.Error())
	}

	return utils.SendSuccessResponse(c, response.Success, comment)
}

// DeleteComment menghapus komentar
func (h *CommentController) DeleteComment(c *fiber.Ctx) error {
	commentID, err := strconv.ParseUint(c.Params("commentID"), 10, 32)
	if err != nil {
		return utils.SendErrorResponse(c, response.BadRequest, "ID komentar tidak valid")
	}

	authorID := c.Locals("userID").(uint)

	err = h.newCommentUseCase.DeleteComment(uint(commentID), authorID)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound("")) {
			return utils.SendErrorResponse(c, response.BadRequest, err.Error())
		}
		if errors.Is(err, utils.ErrForbidden("")) {
			return utils.SendErrorResponse(c, response.BadRequest, err.Error())
		}
		return utils.SendErrorResponse(c, response.ServerError, err.Error())
	}

	return utils.SendSuccessResponse(c, response.Success)
}
