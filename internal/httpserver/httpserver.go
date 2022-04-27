package httpserver

import (
	"github.com/alexkopcak/gophermart/internal/auth"
	"github.com/alexkopcak/gophermart/internal/auth/handlers"
	"github.com/gin-contrib/gzip"

	"github.com/gin-gonic/gin"
)

func NewGinEngine(auc auth.UseCase) *gin.Engine {
	router := gin.Default()

	router.Use(gzipMiddlewareHandle)
	router.Use(gzip.Gzip(gzip.BestSpeed, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)))

	handlers.RegisterHTTPEndpoints(router, auc)
	// authHander := NewAuthHandler(auc)
	// router.POST("/api/user/register", authHander.SignUp)
	// router.POST("/api/user/login", authHander.SignIn)

	// router.Use(AuthMiddlewareHandle(auc)).GET("/api/user/orders", authHander.Test)

	return router
}
