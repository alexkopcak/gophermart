package app

import (
	"github.com/alexkopcak/gophermart/internal/auth"
	authlocalstorage "github.com/alexkopcak/gophermart/internal/auth/repository/localstorage"
	authusecase "github.com/alexkopcak/gophermart/internal/auth/usecase"
	"github.com/alexkopcak/gophermart/internal/order"
	orderlocalstorage "github.com/alexkopcak/gophermart/internal/order/repository/localstorage"
	orderusecase "github.com/alexkopcak/gophermart/internal/order/usecase"
	"github.com/gin-gonic/gin"

	"github.com/alexkopcak/gophermart/internal/config"
	"github.com/alexkopcak/gophermart/internal/httpserver"
	"github.com/rs/zerolog/log"
)

type App struct {
	config *config.Config
	server *gin.Engine

	authUC  auth.UseCase
	orderUC order.UseCase
}

func NewApp(cfg *config.Config) *App {
	userRepo := authlocalstorage.NewUserLocalStorage()
	//userRepo := db.NewPostgresStorage(cfg.DataBaseURI)

	orderRepo := orderlocalstorage.NewOrderLocalStorage()

	return &App{
		config: cfg,
		authUC: authusecase.NewAuthUseCase(userRepo,
			cfg.HashSalt,
			cfg.SigningKey,
			cfg.TokenTTL),
		orderUC: orderusecase.NewOrderUseCase(orderRepo),
	}
}

func (app *App) Run() error {
	log.Debug().Str("package", "app").Str("func", "run").Msg("start")

	app.server = httpserver.NewGinEngine(app.authUC, app.orderUC)

	log.Fatal().Err(app.server.Run(app.config.RunAddress))

	log.Debug().Str("package", "app").Str("func", "run").Msg("exit")
	return nil
}
