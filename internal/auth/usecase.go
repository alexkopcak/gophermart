package auth

import (
	"context"

	"github.com/alexkopcak/gophermart/internal/models"
)

const CtxUserKey = "user"

type UseCase interface {
	SignUp(ctx context.Context, userName string, password string) error
	SignIn(ctx context.Context, userName string, password string) (string, error)
	ParseToken(ctx context.Context, accessToken string) (*models.User, error)
}
