package httpserver

import (
	"github.com/alexkopcak/gophermart/internal/auth"
	authhandlers "github.com/alexkopcak/gophermart/internal/auth/handlers"
	"github.com/alexkopcak/gophermart/internal/order"
	orderhandlers "github.com/alexkopcak/gophermart/internal/order/handlers"
	"github.com/gin-contrib/gzip"

	"github.com/gin-gonic/gin"
)

func NewGinEngine(auc auth.UseCase, ouc order.UseCase, asaddress string) *gin.Engine {
	router := gin.Default()

	router.Use(gzipMiddlewareHandle)
	router.Use(gzip.Gzip(gzip.BestSpeed, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)))

	authhandlers.RegisterHTTPEndpoints(router, auc)

	orderhandlers.RegisterHTTPEndpoints(router, authhandlers.AuthMiddlewareHandle(auc), ouc, asaddress)

	return router
}
