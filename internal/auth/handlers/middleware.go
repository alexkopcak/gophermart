package handlers

import (
	"net/http"

	"github.com/alexkopcak/gophermart/internal/auth"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func AuthMiddlewareHandle(auc auth.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := log.With().Str("package", "handlers").Str("func", "AuthMiddlewareHandle").Logger()

		logger.Debug().Msg("enter")
		defer logger.Debug().Msg("exit")

		token, err := c.Cookie("Authorization")
		logger.Debug().Str("token", token).Msg("get token value")

		if token == "" || err != nil {
			logger.Debug().Err(err).Msg("exit with error")
			c.String(http.StatusUnauthorized, "пользователь не аутентифицирован")
			return
		}

		user, err := auc.ParseToken(c.Request.Context(), token)
		logger.Debug().Str("user", user.UserName).Msg("user from jw token")

		if err != nil || user.UserName == "" {
			logger.Debug().Err(err).Msg("exit with error")
			c.String(http.StatusUnauthorized, "пользователь не аутентифицирован")
			return
		}
		c.Set(auth.CtxUserKey, user.UserName)
		c.Next()
	}
}
