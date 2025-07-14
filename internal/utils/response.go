package utils

import "github.com/gofiber/fiber/v2"

type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
}

func SendSuccessResponse(c *fiber.Ctx, statusCode int, message string, data ...interface{}) error {
	resp := SuccessResponse{
		Success: true,
		Message: message,
	}

	if len(data) > 0 {
		resp.Data = data[0]
	}
	return c.Status(statusCode).JSON(resp)
}

func SendErrorResponse(c *fiber.Ctx, statusCode int, message string, errors ...interface{}) error {
	resp := ErrorResponse{
		Success: false,
		Message: message,
		Errors:  errors,
	}

	if len(errors) > 0 {
		resp.Errors = errors[0]
	}

	return c.Status(statusCode).JSON(resp)
}
