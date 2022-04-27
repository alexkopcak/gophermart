package handlers

import (
	"github.com/alexkopcak/gophermart/internal/auth"
	"github.com/gin-gonic/gin"
)

func RegisterHTTPEndpoints(router *gin.Engine, auc auth.UseCase) {
	handler := NewAuthHandler(auc)

	router.POST("/api/user/register", handler.SignUp)
	router.POST("/api/user/login", handler.SignIn)

}
