package handlers

import (
	"sync"

	"github.com/alexkopcak/gophermart/internal/order"
	"github.com/gin-gonic/gin"
)

func RegisterHTTPEndpoints(wg *sync.WaitGroup, uChannel chan *string, router *gin.Engine, midlleware gin.HandlerFunc, ouc order.UseCase, asAddress string) {
	handler := NewOrderHandler(wg, uChannel, ouc, asAddress)
	handler.UpdateNotFinnalizedOrders()

	routes := router.Use(midlleware)

	routes.POST("/api/user/orders", handler.AddNewOrder)
	routes.GET("/api/user/orders", handler.GetUserOrders)
	routes.GET("/api/user/balance", handler.GetUserBalance)
	routes.POST("/api/user/balance/withdraw", handler.BalanceWithdraw)
	routes.GET("/api/user/withdrawals", handler.Withdrawals)
}
