package app

import (
	"github.com/alexkopcak/gophermart/internal/auth"
	"github.com/alexkopcak/gophermart/internal/auth/repository/localstorage"
	"github.com/alexkopcak/gophermart/internal/auth/usecase"
	"github.com/gin-gonic/gin"

	"github.com/alexkopcak/gophermart/internal/config"
	"github.com/alexkopcak/gophermart/internal/httpserver"
	"github.com/rs/zerolog/log"
)

type App struct {
	config *config.Config
	server *gin.Engine

	authUC auth.UseCase
}

func NewApp(cfg *config.Config) *App {
	userRepo := localstorage.NewUserLocalStorage()

	return &App{
		config: cfg,
		authUC: usecase.NewAuthUseCase(userRepo,
			cfg.HashSalt,
			cfg.SigningKey,
			cfg.TokenTTL),
	}
}

func (app *App) Run() error {
	log.Debug().Str("package", "app").Str("func", "run").Msg("start")

	app.server = httpserver.NewGinEngine(app.authUC)

	log.Fatal().Err(app.server.Run(app.config.RunAddress))

	log.Debug().Str("package", "app").Str("func", "run").Msg("exit")
	return nil
}
