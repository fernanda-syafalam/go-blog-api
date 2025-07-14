// handlers/comment_handler.go
package http

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/usecase"
	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type CommentController struct {
	newCommentUseCase usecase.CommentUseCase
	validator      *validator.Validate
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
		return utils.SendErrorResponse(c, http.StatusBadRequest, "ID postingan tidak valid")
	}

	var req CreateCommentRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	if err := h.validator.Struct(req); err != nil {
		return utils.SendValidatorErrorResponse(c, err)
	}

	authorID := c.Locals("userID").(uint) 

	comment, err := h.newCommentUseCase.CreateComment(req.Content, uint(postID), authorID)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound("")) {
			return utils.SendErrorResponse(c, http.StatusNotFound, err.Error())
		}
		if strings.Contains(err.Error(), "tidak valid") { // Untuk validasi dari service
			return utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		}
		return utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SendSuccessResponse(c, http.StatusCreated, "Komentar berhasil dibuat", comment)
}

func (h *CommentController) GetCommentsByPostID(c *fiber.Ctx) error {
	postID, err := strconv.ParseUint(c.Params("postID"), 10, 32)
	if err != nil {
		return utils.SendErrorResponse(c, http.StatusBadRequest, "ID postingan tidak valid")
	}

	comments, err := h.newCommentUseCase.GetCommentsByPostID(uint(postID))
	if err != nil {
		if errors.Is(err, utils.ErrNotFound("")) {
			return utils.SendErrorResponse(c, http.StatusNotFound, err.Error())
		}
		return utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SendSuccessResponse(c, http.StatusOK, "Komentar berhasil diambil", comments)
}

type UpdateCommentRequest struct {
	Content string `json:"content" validate:"required,min=1,max=500"`
}

func (h *CommentController) UpdateComment(c *fiber.Ctx) error {
	commentID, err := strconv.ParseUint(c.Params("commentID"), 10, 32)
	if err != nil {
		return utils.SendErrorResponse(c, http.StatusBadRequest, "ID komentar tidak valid")
	}

	var req UpdateCommentRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	if err := h.validator.Struct(req); err != nil {
		return utils.SendValidatorErrorResponse(c, err)
	}

	authorID := c.Locals("userID").(uint)

	comment, err := h.newCommentUseCase.UpdateComment(uint(commentID), authorID, req.Content)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound("")) {
			return utils.SendErrorResponse(c, http.StatusNotFound, err.Error())
		}
		if errors.Is(err, utils.ErrForbidden("")) {
			return utils.SendErrorResponse(c, http.StatusForbidden, err.Error())
		}
		if strings.Contains(err.Error(), "tidak valid") { // Untuk validasi dari service
			return utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		}
		return utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SendSuccessResponse(c, http.StatusOK, "Komentar berhasil diperbarui", comment)
}

// DeleteComment menghapus komentar
func (h *CommentController) DeleteComment(c *fiber.Ctx) error {
	commentID, err := strconv.ParseUint(c.Params("commentID"), 10, 32)
	if err != nil {
		return utils.SendErrorResponse(c, http.StatusBadRequest, "ID komentar tidak valid")
	}

	authorID := c.Locals("userID").(uint)

	err = h.newCommentUseCase.DeleteComment(uint(commentID), authorID)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound("")) {
			return utils.SendErrorResponse(c, http.StatusNotFound, err.Error())
		}
		if errors.Is(err, utils.ErrForbidden("")) {
			return utils.SendErrorResponse(c, http.StatusForbidden, err.Error())
		}
		return utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SendSuccessResponse(c, http.StatusNoContent, "Komentar berhasil dihapus")
}