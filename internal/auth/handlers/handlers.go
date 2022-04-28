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
	log.Debug().Str("package", "handlers").Str("func", "SignUp").Msg("start")
	var user models.User
	//defer c.Request.Body.Close()

	if strings.Compare(c.ContentType(), "application/json") != 0 {
		c.String(http.StatusBadRequest, "неверный формат запроса")
		log.Debug().Str("package", "handlers").Str("func", "SignUp").Str("ContentType", c.ContentType()).Msg("exit with error: bad content type")
		return
	}

	var err error
	err = json.NewDecoder(c.Request.Body).Decode(&user)

	if err != nil || user.UserName == "" {
		c.String(http.StatusBadRequest, "неверный формат запроса", err.Error())
		log.Debug().Str("package", "handlers").Str("func", "SignUp").Str("user", user.UserName).Msg("exit with error: empty user name")
		return
	}

	err = h.AuthUseCase.SignUp(c.Request.Context(), user.UserName, user.Password)
	if errors.Is(err, auth.ErrUserAlreadyExsist) {
		c.String(http.StatusConflict, "логин уже занят")
		log.Debug().Str("package", "handlers").Str("func", "SignUp").Msg("exit with error: user already exsist")
		return
	}
	if err != nil {
		c.String(http.StatusInternalServerError, "внутренняя ошибка сервера")
		log.Debug().Str("package", "handlers").Str("func", "SignUp").Msg("exit with error: something went wrong")
		return
	}

	token, err := h.AuthUseCase.SignIn(c.Request.Context(), user.UserName, user.Password)
	if err != nil {
		c.String(http.StatusInternalServerError, "внутренняя ошибка сервера")
		log.Debug().Str("package", "handlers").Str("func", "SignUp").Msg("exit with error")
		return
	}

	c.SetCookie("Authorization", token, 3600, "/", "", false, false)
	c.String(http.StatusOK, "пользователь успешно зарегистрирован и аутентифицирован")
	log.Debug().Str("package", "handlers").Str("func", "SignUp").Msg("user autentithication success")
	log.Debug().Str("package", "handlers").Str("func", "SignUp").Msg("exit")
}

func (h *AuthHandler) SignIn(c *gin.Context) {
	var user models.User
	defer c.Request.Body.Close()

	if strings.Compare(c.ContentType(), "application/json") != 0 {
		c.String(http.StatusBadRequest, "неверный формат запроса")
		return
	}

	var err error
	err = json.NewDecoder(c.Request.Body).Decode(&user)

	if err != nil || user.UserName == "" {
		c.String(http.StatusBadRequest, "неверный формат запроса")
		return
	}

	token, err := h.AuthUseCase.SignIn(c.Request.Context(), user.UserName, user.Password)
	if errors.Is(err, auth.ErrUserNotExsist) {
		c.String(http.StatusUnauthorized, "неверная пара логин/пароль")
		return
	}
	if err != nil {
		c.String(http.StatusInternalServerError, "внутренняя ошибка сервера")
		return
	}

	c.SetCookie("Authorization", token, 3600, "/", "", false, false)
	c.String(http.StatusOK, "пользователь успешно аутентифицирован")
}

func (h *AuthHandler) Test(c *gin.Context) {
	user := c.GetString(auth.CtxUserKey)
	c.String(http.StatusOK, "вызов Test от "+user)
}
