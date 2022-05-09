package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/alexkopcak/gophermart/internal/auth"
	"github.com/alexkopcak/gophermart/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type AuthHandler struct {
	AuthUseCase auth.UseCase
}

func NewAuthHandler(auc auth.UseCase) *AuthHandler {
	return &AuthHandler{
		AuthUseCase: auc,
	}
}

func (h *AuthHandler) SignUp(c *gin.Context) {
	logger := log.With().Str("package", "handlers").Str("function", "SignUp").Logger()
	logger.Debug().Msg("enter")
	defer logger.Debug().Msg("exit")

	var user models.User
	defer c.Request.Body.Close()

	if strings.Compare(c.ContentType(), "application/json") != 0 {
		logger.Debug().Str("ContentType", c.ContentType()).Msg("exit with error: bad content type")
		c.String(http.StatusBadRequest, "неверный формат запроса")
		return
	}

	var err error
	err = json.NewDecoder(c.Request.Body).Decode(&user)

	if err != nil || user.UserName == "" {
		logger.Debug().Str("user", user.UserName).Msg("exit with error: empty user name")
		c.String(http.StatusBadRequest, "неверный формат запроса", err.Error())
		return
	}

	err = h.AuthUseCase.SignUp(c.Request.Context(), user.UserName, user.Password)
	if errors.Is(err, auth.ErrUserAlreadyExsist) {
		logger.Debug().Msg("exit with error: user already exsist")
		c.String(http.StatusConflict, "логин уже занят")
		return
	}
	if err != nil {
		logger.Debug().Err(err).Msg("exit with error: something went wrong")
		c.String(http.StatusInternalServerError, "внутренняя ошибка сервера")
		return
	}

	token, err := h.AuthUseCase.SignIn(c.Request.Context(), user.UserName, user.Password)
	if err != nil {
		logger.Debug().Err(err).Msg("exit with error")
		c.String(http.StatusInternalServerError, "внутренняя ошибка сервера")
		return
	}

	c.SetCookie("Authorization", token, 3600, "/", "", false, false)
	c.String(http.StatusOK, "пользователь успешно зарегистрирован и аутентифицирован")
	logger.Debug().Msg("user autentithication success")
}

func (h *AuthHandler) SignIn(c *gin.Context) {
	logger := log.With().Str("package", "handlers").Str("function", "SignIn").Logger()
	logger.Debug().Msg("enter")
	defer logger.Debug().Msg("exit")

	var user models.User
	defer c.Request.Body.Close()

	if strings.Compare(c.ContentType(), "application/json") != 0 {
		logger.Debug().Str("ContentType", c.ContentType()).Msg("exit with error: bad content type")
		c.String(http.StatusBadRequest, "неверный формат запроса")
		return
	}

	var err error
	err = json.NewDecoder(c.Request.Body).Decode(&user)

	if err != nil || user.UserName == "" {
		logger.Debug().Str("user.UserName", user.UserName).Msg("exit with error: empty user name")
		c.String(http.StatusBadRequest, "неверный формат запроса")
		return
	}

	token, err := h.AuthUseCase.SignIn(c.Request.Context(), user.UserName, user.Password)
	if errors.Is(err, auth.ErrUserNotExsist) {
		logger.Debug().Msg("exit with error: bad user name or password")
		c.String(http.StatusUnauthorized, "неверная пара логин/пароль")
		return
	}
	if err != nil {
		logger.Debug().Str("user.UserName", user.UserName).Err(err).Msg("exit with error")
		c.String(http.StatusInternalServerError, "внутренняя ошибка сервера")
		return
	}

	logger.Debug().Msg("user authorization complete")
	c.SetCookie("Authorization", token, 3600, "/", "", false, false)
	c.String(http.StatusOK, "пользователь успешно аутентифицирован")
}

func (h *AuthHandler) Test(c *gin.Context) {
	user := c.GetString(auth.CtxUserKey)
	c.String(http.StatusOK, "вызов Test от "+user)
}
