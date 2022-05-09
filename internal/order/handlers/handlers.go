package handlers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

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

func NewOrderHandler(ouc order.UseCase, accrualServiceAddress string) *OrderHandler {
	accrualService := integration.NewAccurualService(accrualServiceAddress, ouc)

	return &OrderHandler{
		OrderUseCase:    ouc,
		AccurualService: accrualService,
	}
}

func (h *OrderHandler) AddNewOrder(c *gin.Context) {
	log.Logger = log.With().Str("package", "handlers").Str("function", "AddNewOrder").Logger()
	log.Debug().Msg("enter")
	defer log.Debug().Msg("exit")

	if strings.Compare(c.ContentType(), "text/plain") != 0 {
		log.Debug().Str("ContentType", c.ContentType()).Msg("exit with error: bad content type")
		c.String(http.StatusBadRequest, "неверный формат запроса")
		return
	}

	buff, err := ioutil.ReadAll(c.Request.Body)
	var orderID = string(buff)
	log.Debug().Str("order", orderID).Msg("get order number from request body")

	if err != nil || orderID == "" {
		log.Debug().Err(err).Msg("exit with error")
		c.String(http.StatusUnprocessableEntity, "неверный формат номера заказа")
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		log.Debug().Err(err).Msg("exit with error")
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	err = h.OrderUseCase.AddNewOrder(c.Request.Context(), userID, orderID)
	if err != nil {
		log.Debug().Err(err).Msg("exit with error")
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

	go func() {
		h.AccurualService.UpdateData(orderID)
	}()
	c.String(http.StatusAccepted, "новый номер заказа принят в обработку")
	log.Debug().Msg("new order has accepted")
}

func getUserID(c *gin.Context) (string, error) {
	user, exsists := c.Get(auth.CtxUserKey)
	if !exsists {
		c.String(http.StatusUnauthorized, "пользователь не аутентифицирован")
		return "", order.ErrUserNotAuthtorised
	}
	userID, ok := user.(string)
	if !ok {
		c.String(http.StatusInternalServerError, "не могу получить пользователя")
		return userID, errors.New("не могу получить пользователя")
	}
	return userID, nil
}

func (h *OrderHandler) GetUserOrders(c *gin.Context) {
	log.Logger = log.With().Str("package", "handlers").Str("function", "GetUserOrders").Logger()
	log.Debug().Msg("enter")
	defer log.Debug().Msg("exit")

	userID, err := getUserID(c)
	log.Debug().Str("user", userID).Msg("get user ID")
	if err != nil {
		log.Debug().Err(err).Msg("exit with error")
		c.String(http.StatusInternalServerError, "внутренняя ошибка сервера")
		return
	}

	orders, err := h.OrderUseCase.GetOrders(c.Request.Context(), userID)
	log.Debug().Msg("Get user orders")
	if err != nil {
		log.Debug().Err(err).Msg("exit with error")
		c.String(http.StatusInternalServerError, "внутренняя ошибка сервера")
		return
	}

	if len(orders) == 0 {
		log.Debug().Msg("where are no orders")
		c.JSON(http.StatusNoContent, orders)
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) GetUserBalance(c *gin.Context) {
	log.Logger = log.With().Str("package", "handlers").Str("function", "GetUserBalance").Logger()
	log.Debug().Msg("enter")
	defer log.Debug().Msg("exit")

	userID, err := getUserID(c)
	if err != nil {
		log.Debug().Err(err).Msg("exit with error")
		c.String(http.StatusInternalServerError, "внутренняя ошибка сервера")
		return
	}

	balance, err := h.OrderUseCase.GetBalance(c.Request.Context(), userID)
	log.Debug().Float32("current", balance.Current).Float32("withdrawn", balance.Withdrawn).Msg("get user balance")
	if err != nil {
		log.Debug().Err(err).Msg("exit with error")
		c.String(http.StatusInternalServerError, "внутренняя ошибка сервера")
		return
	}
	c.JSON(http.StatusOK, balance)
}

func (h *OrderHandler) BalanceWithdraw(c *gin.Context) {
	log.Logger = log.With().Str("package", "handlers").Str("function", "BalanceWithdraw").Logger()
	log.Debug().Msg("enter")
	defer log.Debug().Msg("exit")

	if strings.Compare(c.ContentType(), "application/json") != 0 {
		log.Debug().Str("ContentType", c.ContentType()).Msg("exit with error: bad content type")
		c.String(http.StatusBadRequest, "неверный формат запроса")
		return
	}

	var balWithdraw models.BalanceWithdraw
	err := json.NewDecoder(c.Request.Body).Decode(&balWithdraw)

	if err != nil || balWithdraw.OrderID == "" {
		log.Debug().Err(err).Msg("exit with error")
		c.String(http.StatusUnprocessableEntity, "неверный номер заказа")
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		log.Debug().Err(err).Msg("exit with error")
		return
	}
	err = h.OrderUseCase.BalanceWithdraw(c.Request.Context(), userID, &balWithdraw)

	if err != nil {
		log.Debug().Err(err).Msg("exit with error")
		if errors.Is(err, order.ErrOrderBadNumber) {
			c.String(http.StatusUnprocessableEntity, "неверный номер заказа")
			return
		}

		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.String(http.StatusOK, "успешная обработка запроса")
	log.Debug().Msg("query was handled succefuly")
}

func (h *OrderHandler) Withdrawals(c *gin.Context) {
	log.Logger = log.With().Str("package", "handlers").Str("function", "Withdrawals").Logger()
	log.Debug().Msg("enter")
	defer log.Debug().Msg("exit")

	userID, err := getUserID(c)
	if err != nil {
		log.Debug().Err(err).Msg("exit with error")
		return
	}

	withdrawls, err := h.OrderUseCase.Withdrawals(c.Request.Context(), userID)
	if err != nil {
		log.Debug().Err(err).Msg("exit with error")
		c.String(http.StatusInternalServerError, "внутренняя ошибка сервера")
		return
	}
	if len(withdrawls) == 0 {
		c.String(http.StatusNoContent, "нет ни одного списания")
		return
	}
	c.JSON(http.StatusOK, withdrawls)
}
