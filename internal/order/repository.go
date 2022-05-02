package order

import (
	"context"

	"github.com/alexkopcak/gophermart/internal/models"
)

type OrderRepository interface {
	InsertOrder(ctx context.Context, userID string, orderNumber string) error
	GetOrdersListByUserID(ctx context.Context, userID string) ([]*models.Order, error)
	GetBalanceByUserID(ctx context.Context, userID string) (*models.Balance, error)
	WithdrawBalance(ctx context.Context, userID string, bw *models.BalanceWithdraw) error
	Withdrawals(ctx context.Context, userID string) ([]*models.Withdrawals, error)
	UpdateOrder(ctx context.Context, orderNumber string, orderStatus string, orderAccrual int32) error
	GetNotFinnalizedOrdersListByUserID(ctx context.Context, userID string) ([]*models.Order, error)
}
