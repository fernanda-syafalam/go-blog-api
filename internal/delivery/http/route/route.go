package route

import (
	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/delivery/http"
	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/delivery/http/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/knadh/koanf"
)

type RouteConfig struct {
	App                *fiber.App
	Config             *koanf.Koanf
	UserController     *http.UserController
	PostController     *http.PostController
	CategoryController *http.CategoryController
	CommentController  *http.CommentController
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	api := c.App.Group("/api/v1")

	auth := api.Group("/auth")
	auth.Post("/register", c.UserController.Register)
	auth.Post("/login", c.UserController.Login)

	posts := api.Group("/posts")
	posts.Get("/", c.PostController.GetAllPosts)
	posts.Get("/:id", c.PostController.GetPostByID)
	posts.Get("/slug/:slug", c.PostController.GetPostBySlug)

	// comments := api.Group("/comments")
	posts.Get("/:postID/comments", c.CommentController.GetCommentsByPostID)

	categories := api.Group("/categories")
	categories.Get("/", c.CategoryController.GetAllCategories)
	categories.Get("/:id", c.CategoryController.GetCategoryByID)
}

func (c *RouteConfig) SetupAuthRoute() {
	api := c.App.Group("/api/v1")
	api.Use(middleware.JWTMiddleware(c.Config))


	auth := api.Group("/auth")
	auth.Get("/me", c.UserController.GetCurrentUser)

	users := api.Group("/users")
	users.Get("/", c.UserController.GetAllUsers)   
	users.Get("/:id", c.UserController.GetUserByID) 
	users.Put("/:id", c.UserController.UpdateUser)
	users.Delete("/:id", c.UserController.DeleteUser)

	post := api.Group("/posts")
	post.Post("/", c.PostController.CreatePost)
	post.Put("/:id", c.PostController.UpdatePost)
	post.Delete("/:id", c.PostController.DeletePost)

	comments := api.Group("/comments")
	post.Post("/:postID/comments", c.CommentController.CreateComment)
	comments.Put("/:commentID", c.CommentController.UpdateComment)
	comments.Delete("/:commentID", c.CommentController.DeleteComment)

	categories := api.Group("/categories")
	categories.Post("/", c.CategoryController.CreateCategory)
	categories.Put("/:id", c.CategoryController.UpdateCategory)

}
