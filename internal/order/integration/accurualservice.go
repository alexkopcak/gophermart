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
	OrderRepository      order.OrderRepository
}

func NewAccurualService(address string, repo order.OrderRepository) *AccurualService {
	return &AccurualService{
		AccrualSystemAddress: address,
		OrderRepository:      repo,
	}
}

var (
	ErrTooManyRequests = errors.New("превышено количество запросов к сервису")
)

func (as *AccurualService) getOrder(number string) (*Order, error) {
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

func (as *AccurualService) UpdateData(number string) error {
	order, err := as.getOrder(number)
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

	err = as.OrderRepository.UpdateOrder(context.Background(), &item)
	if err != nil {
		return err
	}
	return nil
}
