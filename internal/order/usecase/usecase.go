package usecase

import (
	"context"

	"github.com/alexkopcak/gophermart/internal/models"
	"github.com/alexkopcak/gophermart/internal/order"
)

type OrderUseCase struct {
	orderRepo order.OrderRepository
}

func NewOrderUseCase(orderRepo order.OrderRepository) order.UseCase {
	return &OrderUseCase{
		orderRepo: orderRepo,
	}
}

func (ouc *OrderUseCase) AddNewOrder(ctx context.Context, userID string, orderNumber string) error {
	return ouc.orderRepo.InsertOrder(ctx, userID, orderNumber)
}

func (ouc *OrderUseCase) GetOrders(ctx context.Context, userID string) ([]*models.Order, error) {
	return ouc.orderRepo.GetOrdersListByUserID(ctx, userID)
}

func (ouc *OrderUseCase) GetBalance(ctx context.Context, useerID string) (*models.Balance, error) {
	return ouc.orderRepo.GetBalanceByUserID(ctx, useerID)
}

func (ouc *OrderUseCase) BalanceWithdraw(ctx context.Context, userID string, bw *models.BalanceWithdraw) error {
	return ouc.orderRepo.WithdrawBalance(ctx, userID, bw)
}

func (ouc *OrderUseCase) Withdrawals(ctx context.Context, userID string) ([]*models.Withdrawals, error) {
	return nil, nil
}
