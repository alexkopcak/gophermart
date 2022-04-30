package httpserver

import (
	"github.com/alexkopcak/gophermart/internal/auth"
	authhandlers "github.com/alexkopcak/gophermart/internal/auth/handlers"
	"github.com/alexkopcak/gophermart/internal/order"
	orderhandlers "github.com/alexkopcak/gophermart/internal/order/handlers"
	"github.com/gin-contrib/gzip"

	"github.com/gin-gonic/gin"
)

func NewGinEngine(auc auth.UseCase, ouc order.UseCase) *gin.Engine {
	router := gin.Default()

	router.Use(gzipMiddlewareHandle)
	router.Use(gzip.Gzip(gzip.BestSpeed, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)))

	authhandlers.RegisterHTTPEndpoints(router, auc)

	orderhandlers.RegisterHTTPEndpoints(router, authhandlers.AuthMiddlewareHandle(auc), ouc)
	// authHander := NewAuthHandler(auc)
	// router.POST("/api/user/register", authHander.SignUp)
	// router.POST("/api/user/login", authHander.SignIn)

	// router.Use(AuthMiddlewareHandle(auc)).GET("/api/user/orders", authHander.Test)

	return router
}
