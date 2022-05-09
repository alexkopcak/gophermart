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
	log.Logger = log.With().Str("package", "postgres").Str("function", "MakeMigrations").Logger()

	log.Debug().Msg("make migrations")
	defer log.Debug().Msg("exit")
	db, err := sql.Open("postgres", dbURI)
	log.Fatal().Err(err)

	log.Debug().Msg("set driver")
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	log.Fatal().Err(err)

	log.Debug().Msg("set instance")
	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/auth/repository/postgres/migrations",
		"postgres",
		driver)
	log.Fatal().Err(err)

	if m == nil {
		log.Panic().Msg("can't create migrate instance")
	}
	log.Debug().Msg("m up")
	log.Fatal().Err(m.Up())

	log.Debug().Msg("migrate exit")
}
