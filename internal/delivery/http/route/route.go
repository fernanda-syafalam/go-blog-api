package route

import (
	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/delivery/http"
	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App *fiber.App
	UserController *http.UserController
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
}

func (c *RouteConfig) SetupGuestRoute(){
	c.App.Post("api/users", c.UserController.Register)
	c.App.Post("api/users/login", c.UserController.Login)
}


func (c *RouteConfig) SetupAuthRoute(){
	
}