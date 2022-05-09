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

	log.Logger = log.With().Str("package", "main").Str("function", "main").Logger()
	log.Info().Msg("start program")

	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		gin.SetMode(gin.DebugMode)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		gin.SetMode(gin.ReleaseMode)
	}

	cfg := config.Init()

	log.Debug().Str("package", "app").Str("run address", cfg.RunAddress).Msg("get config")

	app := app.NewApp(cfg)

	log.Fatal().Err(app.Run())
	log.Info().Msg("exit program")
}
