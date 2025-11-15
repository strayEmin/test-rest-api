package handler

import (
	"context"
	"errors"
	"regexp"
	"test-task-rest-api/internal/apperrors"
	"test-task-rest-api/internal/transport/rest/dto"
	"unicode/utf8"

	"github.com/gofiber/fiber/v2"
)

type AuthService interface {
	SignUp(ctx context.Context, name, email, password string) (string, error)
	SignIn(ctx context.Context, email, password string) (string, error)
}

type AuthHandler struct {
	authService AuthService
}

func NewAuthHandler(authService AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) SignUp(c *fiber.Ctx) error {
	ctx := c.Context()
	req := dto.UserSignUpRequest{}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse request body",
		})
	}

	if req.Email == "" || req.Pwd == "" || req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Required fields are missing or empty",
		})
	}

	if utf8.RuneCountInString(req.Name) > 32 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name must be between 32 characters long",
		})
	}

	if utf8.RuneCountInString(req.Pwd) < 8 || utf8.RuneCountInString(req.Pwd) > 64 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password must be between 8 and 60 characters long",
		})
	}

	re := regexp.MustCompile("^[^\\s@]+@[^\\s@]+\\.[^\\s@]+$")
	if !re.MatchString(req.Email) || utf8.RuneCountInString(req.Email) > 255 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid email address",
		})
	}

	token, err := h.authService.SignUp(ctx, req.Name, req.Email, req.Pwd)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserAlreadyExists) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "User already exists",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	resp := dto.UserSignUpResponse{
		Token: token,
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

func (h *AuthHandler) SignIn(c *fiber.Ctx) error {
	ctx := c.Context()
	req := dto.UserSignInRequest{}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse request body",
		})
	}

	if utf8.RuneCountInString(req.Pwd) < 8 || utf8.RuneCountInString(req.Pwd) > 64 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password must be between 8 and 60 characters long",
		})
	}

	re := regexp.MustCompile("^[^\\s@]+@[^\\s@]+\\.[^\\s@]+$")
	if !re.MatchString(req.Email) || utf8.RuneCountInString(req.Email) > 255 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid email address",
		})
	}

	token, err := h.authService.SignIn(ctx, req.Email, req.Pwd)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) || errors.Is(err, apperrors.ErrWrongPassword) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Wrong email or password",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	resp := dto.UserSignInResponse{
		Token: token,
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
