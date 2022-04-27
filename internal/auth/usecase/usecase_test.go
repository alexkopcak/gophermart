package usecase

import (
	"context"
	"testing"

	"github.com/alexkopcak/gophermart/internal/auth/repository/mockstorage"
	"github.com/alexkopcak/gophermart/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	repo := new(mockstorage.UserStorageMock)

	uc := NewAuthUseCase(repo, "salt", "secret", 60)

	username := "user"
	password := "password"

	ctx := context.Background()
	user := &models.User{
		UserName: username,
		Password: "c88e9c67041a74e0357befdff93f87dde0904214",
	}

	// Sign Up
	repo.On("CreateUser", user).Return(nil)
	err := uc.SignUp(ctx, username, password)
	assert.NoError(t, err)

	// Sign In
	repo.On("GetUser", user.UserName).Return(user, nil)
	token, err := uc.SignIn(ctx, username, password)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// verify token
	getUser, err := uc.ParseToken(ctx, token)
	assert.NoError(t, err)
	assert.Equal(t, user, getUser)
}
