package main

import (
	"github.com/alexkopcak/gophermart/internal/app"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	debug := true

	log.Info().Msg("start program")

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Fatal().Err(app.Run())
	log.Info().Msg("exit program")
}
