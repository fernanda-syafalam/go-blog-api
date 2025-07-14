package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var Validate *validator.Validate

func InitValidator() {
	Validate = validator.New()
}

type ValidatorError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func SendValidatorErrorResponse(c *fiber.Ctx, err error) error {
	erorrs := make([]ValidatorError, 0)
	for _, err := range err.(validator.ValidationErrors) {
		erorrs = append(erorrs, ValidatorError{
			Field:   err.Field(),
			Message: GetValidationMessage(err),
		})
	}

	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"success": false,
		"message": "validation error",
		"errors":  erorrs,
	})
}

func GetValidationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fe.Field() + " wajib diisi"
	case "min":
		return fe.Field() + " minimal " + fe.Param() + " karakter"
	case "max":
		return fe.Field() + " maksimal " + fe.Param() + " karakter"
	case "email":
		return fe.Field() + " harus berupa alamat email yang valid"
	case "unique": // Contoh custom tag
		return fe.Field() + " sudah digunakan"
	default:
		return fe.Field() + " tidak valid"
	}
}
