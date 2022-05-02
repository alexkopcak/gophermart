package integration

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/alexkopcak/gophermart/internal/models"
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

func (as *AccurualService) getOrder(number string) (*Order, error) {
	var result Order
	for {
		response, err := http.Get(fmt.Sprintf("%s/api/orders/%s", as.AccrualSystemAddress, number))
		if err != nil {
			log.Info().Err(err)
			continue
		}
		defer response.Body.Close()

		if response.StatusCode == http.StatusInternalServerError {
			return nil, nil
		}

		if response.StatusCode == http.StatusTooManyRequests {
			timeSleepString := response.Header.Get("Retry-After")
			timeSleep, err := strconv.Atoi(timeSleepString)
			log.Debug().Str("Retry-After", timeSleepString).Msg("catch timeout")
			log.Debug().Err(err)
			if err != nil {
				continue
			}
			time.Sleep(time.Duration(timeSleep) * time.Second)
			continue
		}

		if response.StatusCode == http.StatusOK {
			err = json.NewDecoder(response.Body).Decode(&result)
			if err != nil {
				return nil, err
			}
			log.Info().Str("Number", result.Number).Str("Status", result.Status).Float32("Accurual", result.Accrual).Msg("получено")
			if result.Status == models.OrderStatusProcessing {
				as.OrderUseCase.UpdateOrder(context.Background(), result.Number, result.Status, 0)
			}
			if result.Status == models.OrderStatusProcessed ||
				result.Status == models.OrderStatusInvalid {
				return &result, nil
			}
		}
	}
}

func (as *AccurualService) UpdateData(number string) error {
	order, err := as.getOrder(number)
	log.Debug().Err(err)
	if err != nil {
		return err
	}
	err = as.OrderUseCase.UpdateOrder(context.Background(), order.Number, order.Status, int32(order.Accrual*100))
	return err
}
