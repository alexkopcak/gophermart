package httpserver

import (
	"sync"

	"github.com/alexkopcak/gophermart/internal/auth"
	authhandlers "github.com/alexkopcak/gophermart/internal/auth/handlers"
	"github.com/alexkopcak/gophermart/internal/order"
	orderhandlers "github.com/alexkopcak/gophermart/internal/order/handlers"
	"github.com/gin-contrib/gzip"

	"github.com/gin-gonic/gin"
)

func NewGinEngine(wg *sync.WaitGroup, uChannel chan *string, auc auth.UseCase, ouc order.UseCase, asaddress string) *gin.Engine {
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.Use(gzipMiddlewareHandle)
	router.Use(gzip.Gzip(gzip.BestSpeed, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)))

	authhandlers.RegisterHTTPEndpoints(router, auc)

	orderhandlers.RegisterHTTPEndpoints(wg, uChannel, router, authhandlers.AuthMiddlewareHandle(auc), ouc, asaddress)

	return router
}
