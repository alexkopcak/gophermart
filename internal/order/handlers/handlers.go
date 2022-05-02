package handlers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/alexkopcak/gophermart/internal/auth"
	"github.com/alexkopcak/gophermart/internal/models"
	"github.com/alexkopcak/gophermart/internal/order"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/theplant/luhn"
)

type OrderHandler struct {
	OrderUseCase order.UseCase
}

func NewOrderHandler(ouc order.UseCase) *OrderHandler {
	return &OrderHandler{
		OrderUseCase: ouc,
	}
}

func (h *OrderHandler) AddNewOrder(c *gin.Context) {
	log.Debug().Str("package", "handlers").Str("func", "AddNewOrder").Msg("start")

	if strings.Compare(c.ContentType(), "text/plain") != 0 {
		log.Debug().Str("package", "handlers").Str("func", "AddNewOrder").Str("ContentType", c.ContentType()).Msg("exit with error: bad content type")
		c.String(http.StatusBadRequest, "неверный формат запроса")
		return
	}

	buff, err := ioutil.ReadAll(c.Request.Body)
	var orderID = string(buff)

	if err != nil || orderID == "" {
		log.Debug().Str("package", "handlers").Str("func", "AddNewOrder").Str("ContentType", c.ContentType()).Msg("exit with error: bad content type")
		c.String(http.StatusUnprocessableEntity, "неверный формат номера заказа")
		return
	}

	id, err := strconv.Atoi(orderID)

	if err != nil || !luhn.Valid(id) {
		c.String(http.StatusUnprocessableEntity, "неверный формат номера заказа")
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		return
	}

	err = h.OrderUseCase.AddNewOrder(c.Request.Context(), userID, orderID)
	if err != nil {
		if errors.Is(err, order.ErrOrderAlreadyInsertedByOtherUser) {
			c.String(http.StatusConflict, "номер заказа уже был загружен другим пользователем")
			return
		}
		if errors.Is(err, order.ErrOrderAlreadyInsertedByUser) {
			c.String(http.StatusOK, "номер заказа уже был загружен этим пользователем")
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusAccepted, "новый номер заказа принят в обработку")
	log.Debug().Str("package", "handlers").Str("func", "AddNewOrder").Msg("exit")
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
	log.Debug().Str("package", "handlers").Str("func", "GetUserOrders").Msg("start")

	userID, err := getUserID(c)
	if err != nil {
		c.String(http.StatusInternalServerError, "внутренняя ошибка сервера")
		return
	}

	orders, err := h.OrderUseCase.GetOrders(c.Request.Context(), userID)
	if err != nil {
		c.String(http.StatusInternalServerError, "внутренняя ошибка сервера")
		return
	}
	if len(orders) == 0 {
		c.String(http.StatusNoContent, "нет данных для ответа")
		return
	}
	c.JSON(http.StatusOK, orders)

	log.Debug().Str("package", "handlers").Str("func", "GetUserOrders").Msg("exit")
}

func (h *OrderHandler) GetUserBalance(c *gin.Context) {
	log.Debug().Str("package", "handlers").Str("func", "GetUserBalance").Msg("start")

	userID, err := getUserID(c)
	if err != nil {
		return
	}

	balance, err := h.OrderUseCase.GetBalance(c.Request.Context(), userID)
	if err != nil {
		c.String(http.StatusInternalServerError, "внутренняя ошибка сервера")
		return
	}
	c.JSON(http.StatusOK, balance)

	log.Debug().Str("package", "handlers").Str("func", "GetUserBalance").Msg("exit")
}

func (h *OrderHandler) BalanceWithdraw(c *gin.Context) {
	log.Debug().Str("package", "handlers").Str("func", "BalanceWithdraw").Msg("start")

	if strings.Compare(c.ContentType(), "application/json") != 0 {
		log.Debug().Str("package", "handlers").Str("func", "BalanceWithdraw").Str("ContentType", c.ContentType()).Msg("exit with error: bad content type")
		c.String(http.StatusBadRequest, "неверный формат запроса")
		return
	}

	var balWithdraw models.BalanceWithdraw
	err := json.NewDecoder(c.Request.Body).Decode(&balWithdraw)

	if err != nil || balWithdraw.OrderID == "" {
		log.Debug().Str("package", "handlers").Str("func", "BalanceWithdraw").Str("OrderID", balWithdraw.OrderID).Msg("start")
		c.String(http.StatusUnprocessableEntity, "неверный номер заказа")
		return
	}

	id, err := strconv.Atoi(balWithdraw.OrderID)

	if err != nil || !luhn.Valid(id) {
		c.String(http.StatusUnprocessableEntity, "неверный формат номера заказа")
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		return
	}
	err = h.OrderUseCase.BalanceWithdraw(c.Request.Context(), userID, &balWithdraw)

	if errors.Is(err, order.ErrOrderBadNumber) {
		c.String(http.StatusUnprocessableEntity, "неверный номер заказа")
		return
	}

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.String(http.StatusOK, "успешная обработка запроса")
	log.Debug().Str("package", "handlers").Str("func", "BalanceWithdraw").Msg("exit")
}

func (h *OrderHandler) Withdrawals(c *gin.Context) {
	log.Debug().Str("package", "handlers").Str("func", "Withdrawals").Msg("start")
	userID, err := getUserID(c)
	if err != nil {
		return
	}

	withdrawls, err := h.OrderUseCase.Withdrawals(c.Request.Context(), userID)
	if err != nil {
		c.String(http.StatusInternalServerError, "внутренняя ошибка сервера")
		return
	}
	if len(withdrawls) == 0 {
		c.String(http.StatusNoContent, "нет ни одного списания")
		return
	}
	c.JSON(http.StatusOK, withdrawls)

	log.Debug().Str("package", "handlers").Str("func", "Withdrawals").Msg("exit")
}
