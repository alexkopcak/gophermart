package integration

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/alexkopcak/gophermart/internal/models"
	"github.com/alexkopcak/gophermart/internal/order"
)

type Order struct {
	Number  string  `json:"number"`
	Status  string  `json:"staus"`
	Accrual float32 `json:"accrual"`
}

type AccurualService struct {
	AccrualSystemAddress string
	OrderUseCase         order.UseCase
}

func NewAccurualService(address string, usecase order.UseCase) *AccurualService {
	return &AccurualService{
		AccrualSystemAddress: address,
		OrderUseCase:         usecase,
	}
}

var (
	ErrTooManyRequests = errors.New("превышено количество запросов к сервису")
)

func (as *AccurualService) getOrder(ctx context.Context, number string) (*Order, error) {
	var result Order
	response, err := http.Get(fmt.Sprintf("%s/api/orders/%s", as.AccrualSystemAddress, number))

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusTooManyRequests {
		return nil, ErrTooManyRequests
	}

	if response.StatusCode != http.StatusOK {
		return nil, nil
	}

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (as *AccurualService) UpdateData(ctx context.Context, number string) error {
	order, err := as.getOrder(ctx, number)
	if errors.Is(err, ErrTooManyRequests) {
		return err
	}
	if err != nil {
		return err
	}

	item := models.Order{
		Number:  order.Number,
		Status:  order.Status,
		Accrual: order.Accrual,
	}

	err = as.OrderUseCase.UpdateOrder(ctx, &item)
	if err != nil {
		return err
	}
	return nil
}
