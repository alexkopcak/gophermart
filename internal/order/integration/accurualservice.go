package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
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
	WaitGroup            *sync.WaitGroup
	UpdateChannel        chan *string
}

func NewAccurualService(wg *sync.WaitGroup, uc chan *string, address string, usecase order.UseCase) *AccurualService {
	return &AccurualService{
		AccrualSystemAddress: address,
		OrderUseCase:         usecase,
		WaitGroup:            wg,
		UpdateChannel:        uc,
	}
}

func (as *AccurualService) StartUpdateWorker() {
	workerCount := 3

	for i := 0; i < workerCount; i++ {
		as.WaitGroup.Add(1)
		go as.updateWorker(as.WaitGroup, as.UpdateChannel)
	}

}

func (as *AccurualService) updateWorker(wg *sync.WaitGroup, uChan <-chan *string) {
	defer wg.Done()

	for job := range uChan {
		as.UpdateData(*job)
	}
}

func (as *AccurualService) UpdateData(number string) error {
	logger := log.With().Str("package", "integration").Str("function", "UpdateData").Logger()
	logger.Debug().Msg("enter")
	defer logger.Debug().Msg("exit")

	var result Order
	for {
		response, err := http.Get(fmt.Sprintf("%s/api/orders/%s", as.AccrualSystemAddress, number))
		if err != nil {
			logger.Debug().Err(err)
			continue
		}
		defer response.Body.Close()

		if response.StatusCode == http.StatusInternalServerError {
			return nil
		}

		if response.StatusCode == http.StatusTooManyRequests {
			timeSleepString := response.Header.Get("Retry-After")
			timeSleep, err := strconv.Atoi(timeSleepString)
			logger.Debug().Str("Retry-After", timeSleepString).Msg("catch timeout")
			logger.Debug().Err(err).Msg("error message")
			if err != nil {
				continue
			}
			logger.Debug().Int("timeSleep", timeSleep).Msg("wait a some time")
			time.Sleep(time.Duration(timeSleep) * time.Second)
			continue
		}

		if response.StatusCode == http.StatusOK {
			err = json.NewDecoder(response.Body).Decode(&result)
			if err != nil {
				return err
			}

			logger.Debug().Str("response.Status", response.Status).Str("Number", result.Number).Str("Status", result.Status).Float32("Accurual", result.Accrual).Msg("get order")

			/*
				REGISTERED — заказ зарегистрирован, но не начисление не рассчитано;
				INVALID — заказ не принят к расчёту, и вознаграждение не будет начислено;
				PROCESSING — расчёт начисления в процессе;
				PROCESSED — расчёт начисления окончен
			*/

			var status string
			switch result.Status {
			case "REGISTERED":
				status = models.OrderStatusProcessing
			case "INVALID":
				status = models.OrderStatusInvalid
			case "PROCESSING":
				status = models.OrderStatusProcessing
			case "PROCESSED":
				status = models.OrderStatusProcessed
			}

			if status == "" {
				continue
			}

			if result.Status == models.OrderStatusProcessing {
				as.OrderUseCase.UpdateOrder(context.Background(), result.Number, status, 0)
			}
			if result.Status == models.OrderStatusProcessed ||
				result.Status == models.OrderStatusInvalid {
				as.OrderUseCase.UpdateOrder(context.Background(), result.Number, status, int32(result.Accrual*100))
				return nil
			}
		}
	}
}
