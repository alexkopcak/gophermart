package mockstorage

import (
	"context"

	"github.com/alexkopcak/gophermart/internal/models"
	"github.com/stretchr/testify/mock"
)

type UserStorageMock struct {
	mock.Mock
}

func (usm *UserStorageMock) CreateUser(ctx context.Context, user *models.User) error {
	args := usm.Called(user)

	return args.Error(0)
}

func (usm *UserStorageMock) GetUser(ctx context.Context, userName string) (*models.User, error) {
	args := usm.Called(userName)

	return args.Get(0).(*models.User), args.Error(1)
}
