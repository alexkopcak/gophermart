package order

import (
	"context"

	"github.com/alexkopcak/gophermart/internal/models"
)

const CtxUserKey = "user"

type UseCase interface {
	AddNewOrder(ctx context.Context, userID string, orderNumber string) error
	GetOrders(ctx context.Context, userID string) ([]*models.Order, error)
	GetBalance(ctx context.Context, userID string) (*models.Balance, error)
	BalanceWithdraw(ctx context.Context, userID string, bw *models.BalanceWithdraw) error
	Withdrawals(ctx context.Context, userID string) ([]*models.Withdrawals, error)
	UpdateOrder(ctx context.Context, order *models.Order) error
	GetNotFinnalizedOrdersListByUserID(ctx context.Context, userID string) ([]*models.Order, error)
}
