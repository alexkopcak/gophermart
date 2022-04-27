package usecase

import (
	"context"

	"github.com/alexkopcak/gophermart/internal/models"
	"github.com/stretchr/testify/mock"
)

type AuthUseCaseMock struct {
	mock.Mock
}

func (aucm *AuthUseCaseMock) SignUp(ctx context.Context, userName string, password string) error {
	args := aucm.Called(userName, password)

	return args.Error(0)
}

func (aucm *AuthUseCaseMock) SignIn(ctx context.Context, userName string, password string) (string, error) {
	args := aucm.Called(userName, password)

	return args.Get(0).(string), args.Error(1)
}

func (aucm *AuthUseCaseMock) ParseToken(ctx context.Context, accessToken string) (*models.User, error) {
	args := aucm.Called(accessToken)

	return args.Get(0).(*models.User), args.Error(1)
}
