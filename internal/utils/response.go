package utils

import (
	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/common/response"
	"github.com/gofiber/fiber/v2"
)

type SuccessResponse struct {
	ResponseStatus  bool        `json:"responseStatus"`
	ResponseCode    string      `json:"responseCode"`
	ResponseMessage string      `json:"responseMessage"`
	Data            interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	ResponseStatus  bool        `json:"responseStatus"`
	ResponseCode    string      `json:"responseCode"`
	ResponseMessage string      `json:"responseMessage"`
	Errors          interface{} `json:"errors,omitempty"`
}

func SendSuccessResponse(c *fiber.Ctx, code response.Code, data ...interface{}) error {
	resp := SuccessResponse{
		ResponseStatus:  true,
		ResponseCode:    code.GetCode(),
		ResponseMessage: code.GetMessage(),
	}

	if len(data) > 0 {
		resp.Data = data[0]
	}

	return c.Status(code.GetHTTPCode()).JSON(resp)
}

func SendErrorResponse(c *fiber.Ctx, code response.Code, errors ...interface{}) error {
	resp := ErrorResponse{
		ResponseStatus:  false,
		ResponseCode:    code.GetCode(),
		ResponseMessage: code.GetMessage(),
		Errors:          errors,
	}

	if len(errors) > 0 {
		resp.Errors = errors[0]
	}

	return c.Status(code.GetHTTPCode()).JSON(resp)
}

