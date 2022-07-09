package localstorage

import (
	"context"

	"github.com/alexkopcak/gophermart/internal/models"
	"github.com/alexkopcak/gophermart/internal/order"
	"github.com/jackc/pgtype"
)

type OrderItem struct {
	UserID  int32
	Number  string
	Debet   bool
	Status  string
	Accrual int32
	Date    pgtype.Timestamp
}

type OrderLocalStorage struct {
	order []OrderItem
}

func NewOrderLocalStorage() order.OrderRepository {
	return &OrderLocalStorage{
		order: make([]OrderItem, 0),
	}
}

func (ols *OrderLocalStorage) InsertOrder(ctx context.Context, userID int32, orderNumber string) error {
	orderItem, _ := ols.GetOrderByOrderUID(ctx, userID, orderNumber)
	if orderItem != nil {
		if orderItem.UserName == userID {
			return order.ErrOrderAlreadyInsertedByUser
		} else {
			return order.ErrOrderAlreadyInsertedByOtherUser
		}

	}

	item := OrderItem{
		UserID:  userID,
		Number:  orderNumber,
		Debet:   true,
		Status:  models.OrderStatusNew,
		Accrual: 0,
		Date:    pgtype.Timestamp{},
	}
	ols.order = append(ols.order, item)
	return nil
}

func (ols *OrderLocalStorage) GetOrdersListByUserID(ctx context.Context, userID int32) ([]models.Order, error) {
	result := make([]models.Order, 0)
	for _, item := range ols.order {
		if item.UserID == userID {
			resultItem := models.Order{
				UserName: item.UserID,
				Number:   item.Number,
				Status:   item.Status,
				Accrual:  float32(item.Accrual) / 100,
				Uploaded: item.Date,
			}
			result = append(result, resultItem)
		}
	}
	return result, nil
}

func (ols *OrderLocalStorage) GetBalanceByUserID(ctx context.Context, userID int32) (*models.Balance, error) {
	var result = new(models.Balance)
	for _, item := range ols.order {
		if item.UserID == userID {
			if item.Debet {
				result.Current = result.Current + float32(item.Accrual)
			} else {
				result.Withdrawn = result.Withdrawn + float32(item.Accrual)
			}
		}
	}
	return result, nil
}

func (ols *OrderLocalStorage) GetOrderByOrderUID(ctx context.Context, userID int32, orderNumber string) (*models.Order, error) {
	for _, item := range ols.order {
		if item.Number == orderNumber && item.Debet && item.UserID == userID {
			return &models.Order{
				UserName: item.UserID,
				Number:   item.Number,
				Status:   item.Status,
				Accrual:  float32(item.Accrual) / 100,
				Uploaded: item.Date,
			}, nil
		}
	}
	return nil, nil
}

func (ols *OrderLocalStorage) WithdrawBalance(ctx context.Context, userID int32, bw *models.BalanceWithdraw) error {
	balance, err := ols.GetBalanceByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if balance == nil {
		return order.ErrNotEnougthBalance
	}

	currentBalance := balance.Current - balance.Withdrawn
	if currentBalance < bw.Sum {
		return order.ErrNotEnougthBalance
	}

	item := OrderItem{
		UserID:  userID,
		Number:  bw.OrderID,
		Debet:   false,
		Status:  models.OrderStatusWithDrawn,
		Accrual: int32(bw.Sum),
		Date:    pgtype.Timestamp{},
	}

	ols.order = append(ols.order, item)

	return nil
}

func (ols *OrderLocalStorage) Withdrawals(ctx context.Context, userID int32) ([]*models.Withdrawals, error) {
	result := make([]*models.Withdrawals, 0)
	for _, item := range ols.order {
		if item.UserID == userID && !item.Debet {
			resultItem := &models.Withdrawals{
				OrderID:     item.Number,
				Sum:         float32(item.Accrual) / 100,
				ProcessedAt: item.Date.Time,
			}
			result = append(result, resultItem)
		}
	}
	return result, nil
}

func (ols *OrderLocalStorage) UpdateOrder(ctx context.Context, orderNumber string, orderStatus string, orderAccrual int32) error {
	for id, item := range ols.order {
		if item.Number == orderNumber && item.Debet {
			ols.order[id].Status = orderStatus
			ols.order[id].Accrual = orderAccrual
			return nil
		}
	}

	return nil
}

func (ols *OrderLocalStorage) GetNotFinnalizedOrdersListByUserID(ctx context.Context, userID int32) ([]*models.Order, error) {
	result := make([]*models.Order, 0)
	for _, item := range ols.order {
		if item.UserID == userID && (item.Status == models.OrderStatusNew || item.Status == models.OrderStatusProcessing) {
			resultItem := &models.Order{
				UserName: item.UserID,
				Number:   item.Number,
				Status:   item.Status,
				Accrual:  float32(item.Accrual) / 100,
				Uploaded: item.Date,
			}
			result = append(result, resultItem)
		}
	}
	return result, nil

}

func (ols *OrderLocalStorage) GetNotFinnalizedOrdersList(ctx context.Context) ([]*models.Order, error) {
	result := make([]*models.Order, 0)
	for _, item := range ols.order {
		if item.Status == models.OrderStatusNew || item.Status == models.OrderStatusProcessing {
			resultItem := &models.Order{
				UserName: item.UserID,
				Number:   item.Number,
				Status:   item.Status,
				Accrual:  float32(item.Accrual) / 100,
				Uploaded: item.Date,
			}
			result = append(result, resultItem)
		}
	}
	return result, nil

}
