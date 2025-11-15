package handler

import (
	"context"
	"test-task-rest-api/internal/domain"
)

type UserService interface {
	CreateUser(ctx context.Context, name, email, password string) (domain.User, error)
}

type UserHandler struct {
	userService UserService
}

func NewUserHandler(userService UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}
