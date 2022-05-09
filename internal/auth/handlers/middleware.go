package handlers

import (
	"net/http"

	"github.com/alexkopcak/gophermart/internal/auth"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func AuthMiddlewareHandle(auc auth.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Logger = log.With().Str("package", "handlers").Str("func", "AuthMiddlewareHandle").Logger()

		log.Debug().Msg("enter")
		defer log.Debug().Msg("exit")

		token, err := c.Cookie("Authorization")
		log.Debug().Str("token", token).Msg("get token value")

		if token == "" || err != nil {
			log.Debug().Err(err).Msg("exit with error")
			c.String(http.StatusUnauthorized, "пользователь не аутентифицирован")
			return
		}

		user, err := auc.ParseToken(c.Request.Context(), token)
		log.Debug().Str("user", user.UserName).Msg("user from jw token")

		if err != nil || user.UserName == "" {
			log.Debug().Err(err).Msg("exit with error")
			c.String(http.StatusUnauthorized, "пользователь не аутентифицирован")
			return
		}
		c.Set(auth.CtxUserKey, user.UserName)
		c.Next()
	}
}
