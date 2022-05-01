package handlers

import (
	"github.com/alexkopcak/gophermart/internal/order"
	"github.com/gin-gonic/gin"
)

func RegisterHTTPEndpoints(router *gin.Engine, midlleware gin.HandlerFunc, ouc order.UseCase) {
	handler := NewOrderHandler(ouc)

	routes := router.Use(midlleware).Use(AccuralServiceHandler(ouc))

	routes.POST("/api/user/orders", handler.AddNewOrder)
	routes.GET("/api/user/orders", handler.GetUserOrders)
	routes.GET("/api/user/balance", handler.GetUserBalance)
	routes.POST("/api/user/balance/withdraw", handler.BalanceWithdraw)
	routes.GET("/api/user/withdrawals", handler.Withdrawals)
}
