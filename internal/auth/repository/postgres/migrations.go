package postgres

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func MakeMigrations(dbURI string) {
	logger := log.With().Str("package", "postgres").Str("function", "MakeMigrations").Logger()

	logger.Debug().Msg("make migrations")
	defer logger.Debug().Msg("exit")
	db, err := sql.Open("postgres", dbURI)
	logger.Fatal().Err(err)

	logger.Debug().Msg("set driver")
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	logger.Fatal().Err(err)

	logger.Debug().Msg("set instance")
	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/auth/repository/postgres/migrations",
		"postgres",
		driver)
	logger.Fatal().Err(err)

	if m == nil {
		logger.Panic().Msg("can't create migrate instance")
	}
	logger.Debug().Msg("m up")
	logger.Fatal().Err(m.Up())

	logger.Debug().Msg("migrate exit")
}
