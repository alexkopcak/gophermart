package main

import (
	"github.com/alexkopcak/gophermart/internal/app"
	"github.com/alexkopcak/gophermart/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	debug := true

	logger := log.With().Str("package", "main").Str("function", "main").Logger()
	logger.Info().Msg("start program")
	defer logger.Info().Msg("exit program")

	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		gin.SetMode(gin.DebugMode)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		gin.SetMode(gin.ReleaseMode)
	}

	cfg := config.Init()

	logger.Debug().Str("run address", cfg.RunAddress).Msg("get config")

	app := app.NewApp(cfg)

	logger.Fatal().Err(app.Run())
}
