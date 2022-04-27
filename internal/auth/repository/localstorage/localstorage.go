package localstorage

import (
	"context"
	"sync"

	"github.com/alexkopcak/gophermart/internal/auth"

	"github.com/alexkopcak/gophermart/internal/models"
)

type UserLocalStrage struct {
	users map[string]*models.User
	mutex *sync.Mutex
}

func NewUserLocalStorage() auth.UserRepository {
	return &UserLocalStrage{
		users: make(map[string]*models.User),
		mutex: new(sync.Mutex),
	}
}

func (uls *UserLocalStrage) CreateUser(ctx context.Context, user *models.User) error {
	uls.mutex.Lock()
	defer uls.mutex.Unlock()

	if _, exsist := uls.users[user.UserName]; exsist {
		return auth.ErrUserAlreadyExsist
	}

	uls.users[user.UserName] = user
	return nil
}

func (uls *UserLocalStrage) GetUser(ctx context.Context, userName string) (*models.User, error) {
	uls.mutex.Lock()
	defer uls.mutex.Unlock()

	if _, exsist := uls.users[userName]; !exsist {
		return nil, auth.ErrUserNotExsist
	}
	return uls.users[userName], nil
}
