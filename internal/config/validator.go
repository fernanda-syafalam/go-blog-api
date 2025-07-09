package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/knadh/koanf"
)

func NewValidator(k *koanf.Koanf) *validator.Validate {
	return validator.New()
}
