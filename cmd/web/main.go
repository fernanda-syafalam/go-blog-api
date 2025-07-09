package main

import (
	"context"
	"fmt"
	"log"

	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/config"
	"github.com/gofiber/contrib/otelfiber/v2"
)

func main() {

	config.Load()
	k := config.Get()
	tp := config.InitTracer("backend-monitoring-notification", k)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	log := config.NewLogger(k)
	db := config.NewDatabase(k, log)
	validate := config.NewValidator(k)
	app := config.NewFiber(k)
	redis := config.NewRedis(k)

	config.Boostrap(&config.BootstrapConfig{
		DB:       db,
		App:      app,
		Log:      log,
		Validate: validate,
		Config:   k,
		Redis:    redis,
	})
	app.Use(otelfiber.Middleware())

	webPort := k.String("web.port")
	err := app.Listen(fmt.Sprintf(":%s", webPort))
	if err != nil {
		log.Err(err).Msg("Error starting web server")
	}
}
