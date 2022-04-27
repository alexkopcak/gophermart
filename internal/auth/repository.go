package auth

import (
	"context"

	"github.com/alexkopcak/gophermart/internal/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUser(ctx context.Context, userName string) (*models.User, error)
}
