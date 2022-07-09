package localstorage

import (
	"context"
	"testing"

	"github.com/alexkopcak/gophermart/internal/auth"
	"github.com/alexkopcak/gophermart/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	storage := NewUserLocalStorage()

	user := &models.User{
		UserName: "user",
		Password: "password",
	}

	err := storage.CreateUser(context.Background(), user)

	assert.NoError(t, err)

	var getUser *models.User
	getUser, err = storage.GetUser(context.Background(), "user")
	assert.NoError(t, err)
	assert.Equal(t, user, getUser)

	_, err = storage.GetUser(context.Background(), "unknown user")
	assert.Error(t, err)
	assert.Equal(t, err, auth.ErrUserNotExsist)
}
