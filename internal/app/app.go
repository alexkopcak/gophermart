package app

import (
	"sync"

	"github.com/alexkopcak/gophermart/internal/auth"
	authdb "github.com/alexkopcak/gophermart/internal/auth/repository/postgres"
	authusecase "github.com/alexkopcak/gophermart/internal/auth/usecase"
	"github.com/alexkopcak/gophermart/internal/order"

	orderdb "github.com/alexkopcak/gophermart/internal/order/repository/postgres"
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
	logger := log.With().Str("package", "app").Str("function", "NewApp").Logger()
	logger.Debug().Msg("enter")
	defer logger.Debug().Msg("exit")

	//userRepo := authlocalstorage.NewUserLocalStorage()
	userRepo := authdb.NewPostgresStorage(cfg.DataBaseURI)

	//orderRepo := orderlocalstorage.NewOrderLocalStorage()
	orderRepo := orderdb.NewOrderPostgresStorage(cfg.DataBaseURI)

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
	logger := log.With().Str("package", "app").Str("func", "run").Logger()

	wg := &sync.WaitGroup{}
	uChannel := make(chan *string)

	logger.Debug().Msg("enter")
	defer logger.Debug().Msg("exit")

	logger.Debug().Msg("create new gin engine object")
	app.server = httpserver.NewGinEngine(wg, uChannel, app.authUC, app.orderUC, app.config.AccrualSystemAddress)

	logger.Fatal().Err(app.server.Run(app.config.RunAddress))

	close(uChannel)
	wg.Wait()
	return nil
}
