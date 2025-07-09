package config

import (
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/knadh/koanf"
)

func NewFiber(config *koanf.Koanf) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      config.String("app.name"),
		ErrorHandler: NewErrorHandler(),
		Prefork:      config.Bool("web.prefork"),
		JSONEncoder:  sonic.Marshal,
		JSONDecoder:  sonic.Unmarshal,
	})

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))
	return app
}

func NewErrorHandler() fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}

		return ctx.Status(code).JSON(fiber.Map{
			"code":    code,
			"message": err.Error(),
		})
	}
}
