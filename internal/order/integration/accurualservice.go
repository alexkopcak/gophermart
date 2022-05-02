package integration

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/alexkopcak/gophermart/internal/order"
	"github.com/rs/zerolog/log"
)

type Order struct {
	Number  string  `json:"order"`
	Status  string  `json:"status"`
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
	var response *http.Response
	var err error
	var statusCode int
	for ok := true; ok; ok = (statusCode == http.StatusOK) {
		response, err := http.Get(fmt.Sprintf("%s/api/orders/%s", as.AccrualSystemAddress, number))
		statusCode = response.StatusCode
		if err != nil {
			return nil, err
		}
		defer response.Body.Close()

		if response.StatusCode == http.StatusTooManyRequests {
			timeSleepString := response.Header.Get("Retry-After")
			timeSleep, err := strconv.Atoi(timeSleepString)
			log.Debug().Err(err)
			if err != nil {
				return nil, nil
			}
			time.Sleep(time.Duration(timeSleep) * time.Second)
		}

		if response.StatusCode != http.StatusOK {
			return nil, nil
		}
	}

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (as *AccurualService) UpdateData(ctx context.Context, number string) error {
	order, err := as.getOrder(ctx, number)
	log.Debug().Err(err)
	if errors.Is(err, ErrTooManyRequests) {

		return err
	}
	if err != nil {
		return err
	}

	// log.Debug().Str("order number", order.Number).
	// 	Str("order status", order.Status).
	// 	Float32("order accural", order.Accrual).
	// 	Msg("Get from accural data")

	err = as.OrderUseCase.UpdateOrder(ctx, order.Number, order.Status, int32(order.Accrual*100))
	return err
}
