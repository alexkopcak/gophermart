package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/alexkopcak/gophermart/internal/auth"
	"github.com/alexkopcak/gophermart/internal/models"
	"github.com/alexkopcak/gophermart/internal/order"
	"github.com/alexkopcak/gophermart/internal/order/integration"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type OrderHandler struct {
	OrderUseCase    order.UseCase
	AccurualService *integration.AccurualService
}

func NewOrderHandler(wg *sync.WaitGroup, uc chan *string, ouc order.UseCase, accrualServiceAddress string) *OrderHandler {
	accrualService := integration.NewAccurualService(wg, uc, accrualServiceAddress, ouc)
	accrualService.StartUpdateWorker()

	return &OrderHandler{
		OrderUseCase:    ouc,
		AccurualService: accrualService,
	}
}

type orderItem struct {
	models.Order
	Uploaded string `json:"uploaded_at"`
}

func (h *OrderHandler) UpdateNotFinnalizedOrders() {
	logger := log.With().Str("package", "handlers").Str("function", "UpdateUnhandledOrders").Logger()
	logger.Debug().Msg("enter")
	defer logger.Debug().Msg("exit")

	logger.Debug().Msg("Get")
	orders, err := h.OrderUseCase.GetNotFinnalizedOrdersList(context.Background())

	if err != nil {
		logger.Debug().Err(err).Msg("exit with error")
	}
	logger.Debug().Int("len(orders)", len(orders)).Msg("Not finnalized orders count")
	if len(orders) == 0 {
		return
	}

	for _, item := range orders {
		h.AccurualService.UpdateChannel <- &item.Number
	}
}

func (h *OrderHandler) AddNewOrder(c *gin.Context) {
	logger := log.With().Str("package", "handlers").Str("function", "AddNewOrder").Logger()
	logger.Debug().Msg("enter")
	defer logger.Debug().Msg("exit")

	if strings.Compare(c.ContentType(), "text/plain") != 0 {
		logger.Debug().Str("ContentType", c.ContentType()).Msg("exit with error: bad content type")
		c.String(http.StatusBadRequest, "неверный формат запроса")
		return
	}

	buff, err := ioutil.ReadAll(c.Request.Body)
	var orderID = string(buff)
	logger.Debug().Str("order", orderID).Msg("get order number from request body")

	if err != nil || orderID == "" {
		logger.Debug().Err(err).Msg("exit with error")
		c.String(http.StatusUnprocessableEntity, "неверный формат номера заказа")
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		logger.Debug().Err(err).Msg("exit with error")
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	err = h.OrderUseCase.AddNewOrder(c.Request.Context(), userID, orderID)
	if err != nil {
		logger.Debug().Err(err).Msg("exit with error")
		if errors.Is(err, order.ErrOrderAlreadyInsertedByOtherUser) {
			c.String(http.StatusConflict, "номер заказа уже был загружен другим пользователем")
			return
		}
		if errors.Is(err, order.ErrOrderAlreadyInsertedByUser) {
			c.String(http.StatusOK, "номер заказа уже был загружен этим пользователем")
			return
		}
		if errors.Is(err, order.ErrOrderBadNumber) {
			c.String(http.StatusUnprocessableEntity, "неверный формат номера заказа")
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	logger.Debug().Str("orderID", orderID).Msg("orderID sent to accurual service")
	h.AccurualService.UpdateChannel <- &orderID

	c.String(http.StatusAccepted, "новый номер заказа принят в обработку")
	logger.Debug().Msg("new order has accepted")
}

func getUserID(c *gin.Context) (int32, error) {
	user, exsists := c.Get(auth.CtxUserKey)
	if !exsists {
		c.String(http.StatusUnauthorized, "пользователь не аутентифицирован")
		return 0, order.ErrUserNotAuthtorised
	}
	userID, ok := user.(int32)
	if !ok {
		c.String(http.StatusInternalServerError, "не могу получить пользователя")
		return userID, errors.New("не могу получить пользователя")
	}
	return userID, nil
}

func (h *OrderHandler) GetUserOrders(c *gin.Context) {
	logger := log.With().Str("package", "handlers").Str("function", "GetUserOrders").Logger()
	logger.Debug().Msg("enter")
	defer logger.Debug().Msg("exit")

	userID, err := getUserID(c)
	logger.Debug().Int32("user", userID).Msg("get user ID")
	if err != nil {
		logger.Debug().Err(err).Msg("exit with error")
		c.String(http.StatusInternalServerError, "внутренняя ошибка сервера")
		return
	}

	orders, err := h.OrderUseCase.GetOrders(c.Request.Context(), userID)
	logger.Debug().Msg("Get user orders")
	if err != nil {
		logger.Debug().Err(err).Msg("exit with error")
		c.String(http.StatusInternalServerError, "внутренняя ошибка сервера")
		return
	}

	if len(orders) == 0 {
		logger.Debug().Msg("where are no orders")
		c.JSON(http.StatusNoContent, orders)
		return
	}

	var result []orderItem
	for _, item := range orders {
		var resultItem orderItem
		resultItem.Number = item.Number
		resultItem.Status = item.Status
		resultItem.Accrual = item.Accrual
		resultItem.Uploaded = item.Uploaded.Time.Format(time.RFC3339)

		result = append(result, resultItem)
	}

	c.JSON(http.StatusOK, result)
}

func (h *OrderHandler) GetUserBalance(c *gin.Context) {
	logger := log.With().Str("package", "handlers").Str("function", "GetUserBalance").Logger()
	logger.Debug().Msg("enter")
	defer logger.Debug().Msg("exit")

	userID, err := getUserID(c)
	if err != nil {
		logger.Debug().Err(err).Msg("exit with error")
		c.String(http.StatusInternalServerError, "внутренняя ошибка сервера")
		return
	}

	balance, err := h.OrderUseCase.GetBalance(c.Request.Context(), userID)
	logger.Debug().Float32("current", balance.Current).Float32("withdrawn", balance.Withdrawn).Msg("get user balance")
	if err != nil {
		logger.Debug().Err(err).Msg("exit with error")
		c.String(http.StatusInternalServerError, "внутренняя ошибка сервера")
		return
	}
	c.JSON(http.StatusOK, balance)
}

func (h *OrderHandler) BalanceWithdraw(c *gin.Context) {
	logger := log.With().Str("package", "handlers").Str("function", "BalanceWithdraw").Logger()
	logger.Debug().Msg("enter")
	defer logger.Debug().Msg("exit")

	if strings.Compare(c.ContentType(), "application/json") != 0 {
		logger.Debug().Str("ContentType", c.ContentType()).Msg("exit with error: bad content type")
		c.String(http.StatusBadRequest, "неверный формат запроса")
		return
	}

	var balWithdraw models.BalanceWithdraw
	err := json.NewDecoder(c.Request.Body).Decode(&balWithdraw)

	if err != nil || balWithdraw.OrderID == "" {
		logger.Debug().Err(err).Msg("exit with error")
		c.String(http.StatusUnprocessableEntity, "неверный номер заказа")
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		logger.Debug().Err(err).Msg("exit with error")
		return
	}
	err = h.OrderUseCase.BalanceWithdraw(c.Request.Context(), userID, &balWithdraw)

	if err != nil {
		logger.Debug().Err(err).Msg("exit with error")
		if errors.Is(err, order.ErrOrderBadNumber) {
			c.String(http.StatusUnprocessableEntity, "неверный номер заказа")
			return
		}

		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.String(http.StatusOK, "успешная обработка запроса")
	logger.Debug().Msg("query was handled succefuly")
}

func (h *OrderHandler) Withdrawals(c *gin.Context) {
	logger := log.With().Str("package", "handlers").Str("function", "Withdrawals").Logger()
	logger.Debug().Msg("enter")
	defer logger.Debug().Msg("exit")

	userID, err := getUserID(c)
	if err != nil {
		logger.Debug().Err(err).Msg("exit with error")
		return
	}

	withdrawls, err := h.OrderUseCase.Withdrawals(c.Request.Context(), userID)
	if err != nil {
		logger.Debug().Err(err).Msg("exit with error")
		c.String(http.StatusInternalServerError, "внутренняя ошибка сервера")
		return
	}
	if len(withdrawls) == 0 {
		c.String(http.StatusNoContent, "нет ни одного списания")
		return
	}
	c.JSON(http.StatusOK, withdrawls)
}
