package usecase

import (
	"context"
	"strconv"

	"github.com/alexkopcak/gophermart/internal/models"
	"github.com/alexkopcak/gophermart/internal/order"
	"github.com/theplant/luhn"
)

type OrderUseCase struct {
	orderRepo order.OrderRepository
}

func NewOrderUseCase(orderRepo order.OrderRepository) order.UseCase {
	return &OrderUseCase{
		orderRepo: orderRepo,
	}
}

func checkOrderID(orderID string) error {
	id, err := strconv.Atoi(orderID)
	if err != nil {
		return err
	}

	if !luhn.Valid(id) {
		return order.ErrOrderBadNumber
	}
	return nil
}

func (ouc *OrderUseCase) AddNewOrder(ctx context.Context, userID int32, orderNumber string) error {
	err := checkOrderID(orderNumber)
	if err != nil {
		return err
	}
	return ouc.orderRepo.InsertOrder(ctx, userID, orderNumber)
}

func (ouc *OrderUseCase) GetOrders(ctx context.Context, userID int32) ([]models.Order, error) {
	return ouc.orderRepo.GetOrdersListByUserID(ctx, userID)
}

func (ouc *OrderUseCase) GetBalance(ctx context.Context, useerID int32) (*models.Balance, error) {
	return ouc.orderRepo.GetBalanceByUserID(ctx, useerID)
}

func (ouc *OrderUseCase) BalanceWithdraw(ctx context.Context, userID int32, bw *models.BalanceWithdraw) error {
	err := checkOrderID(bw.OrderID)
	if err != nil {
		return err
	}

	return ouc.orderRepo.WithdrawBalance(ctx, userID, bw)
}

func (ouc *OrderUseCase) Withdrawals(ctx context.Context, userID int32) ([]*models.Withdrawals, error) {
	return ouc.orderRepo.Withdrawals(ctx, userID)
}

func (ouc *OrderUseCase) UpdateOrder(ctx context.Context, orderNumber string, orderStatus string, orderAccrual int32) error {
	return ouc.orderRepo.UpdateOrder(ctx, orderNumber, orderStatus, orderAccrual)
}

func (ouc *OrderUseCase) GetNotFinnalizedOrdersListByUserID(ctx context.Context, userID int32) ([]*models.Order, error) {
	return ouc.orderRepo.GetNotFinnalizedOrdersListByUserID(ctx, userID)
}

func (ouc *OrderUseCase) GetNotFinnalizedOrdersList(ctx context.Context) ([]*models.Order, error) {
	return ouc.orderRepo.GetNotFinnalizedOrdersList(ctx)
}
