package handlers

import (
	"github.com/alexkopcak/gophermart/internal/order"
	"github.com/gin-gonic/gin"
)

func RegisterHTTPEndpoints(router *gin.Engine, midlleware gin.HandlerFunc, auc order.UseCase) {
	handler := NewOrderHandler(auc)

	router.Use(midlleware).POST("/api/user/orders", handler.AddNewOrder)
	router.Use(midlleware).GET("/api/user/orders", handler.GetUserOrders)
	router.Use(midlleware).GET("/api/user/balance", handler.GetUserBalance)
	router.Use(midlleware).POST("/api/user/balance/withdraw", handler.BalanceWithdraw)
	router.Use(midlleware).GET("/api/user/withdrawals", handler.Withdrawals)
}
