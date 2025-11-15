package handler

import (
	"context"
	"test-task-rest-api/internal/config"
	"time"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	httpLog "github.com/gofiber/fiber/v2/middleware/logger"
)

type Deps struct {
	AuthService AuthService
}

type Handler struct {
	config *config.Config
	app    *fiber.App

	AuthHandler *AuthHandler
}

func NewHandler(config *config.Config, deps Deps) *Handler {
	return &Handler{
		config:      config,
		AuthHandler: NewAuthHandler(deps.AuthService),
	}
}

func (h *Handler) newJwtCheckingMiddleware() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: h.config.SecretKey},
		ContextKey: "jwt",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})
}

func (h *Handler) initRoutes() {
	api := h.app.Group("/api")
	v1 := api.Group("/v1")

	auth := v1.Group("/auth")
	auth.Post("signup", h.AuthHandler.SignUp)
	auth.Post("signin", h.AuthHandler.SignIn)

	v1.Use(h.newJwtCheckingMiddleware())

	transactions := v1.Group("/transactions")
	transactions.Get("", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"error": "transactions.Get not implemented",
		})
	})
	transactions.Patch("/:id", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"id": c.Params("id"),
		})
	})
	transactions.Post("", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"error": "transactions.Post not implemented",
		})
	})
}

func (h *Handler) Init(ctx context.Context) *fiber.App {
	h.app = fiber.New(
		fiber.Config{
			ReadTimeout:  h.config.Timeout,
			WriteTimeout: h.config.Timeout,
			IdleTimeout:  h.config.IdleTimeout,
		})

	h.app.Use(httpLog.New())

	h.app.Use(
		func(c *fiber.Ctx) error {
			ctx, cancel := context.WithTimeout(c.UserContext(), h.config.HTTPServer.Timeout-time.Millisecond*500)
			defer cancel()
			c.SetUserContext(ctx)
			return c.Next()
		})

	h.initRoutes()

	return h.app
}
