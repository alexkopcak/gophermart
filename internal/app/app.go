package app

import "github.com/rs/zerolog/log"

func Run() error {
	log.Debug().Str("package", "app").Str("func", "run").Msg("start")

	log.Debug().Str("package", "app").Str("func", "run").Msg("exit")
	return nil
}
