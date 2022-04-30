package handlers

import (
	"net/http"

	"github.com/alexkopcak/gophermart/internal/auth"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func AuthMiddlewareHandle(auc auth.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Debug().Str("package", "httpserver").Str("func", "authmiddlewarehandle").Msg("enter")

		token, err := c.Cookie("Authorization")
		log.Debug().Str("package", "httpserver").Str("func", "authmiddlewarehandle").Str("token", token).Msg("get token value")

		if token == "" || err != nil {
			log.Debug().Str("package", "httpserver").Str("func", "authmiddlewarehandle").Msg("exit with error")
			c.String(http.StatusUnauthorized, "пользователь не аутентифицирован")
			c.Abort()
			return
		}

		user, err := auc.ParseToken(c.Request.Context(), token)
		if user != nil {
			log.Debug().Str("package", "httpserver").Str("func", "authmiddlewarehandle").Str("user", user.UserName).Str("pass", user.Password).Msg("")
		}
		if err != nil {
			log.Debug().Str("package", "httpserver").Str("func", "authmiddlewarehandle").Str("error", err.Error()).Msg("")
		}

		if err != nil || user.UserName == "" {
			log.Debug().Str("package", "httpserver").Str("func", "authmiddlewarehandle").Msg("exit with error")
			c.String(http.StatusUnauthorized, "пользователь не аутентифицирован")
			c.Abort()
			return
		}
		log.Debug().Str("package", "httpserver").Str("func", "authmiddlewarehandle").Msg("exit")

		c.Set(auth.CtxUserKey, user.UserName)
		c.Next()
	}
}
