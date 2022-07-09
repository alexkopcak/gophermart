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
	log.Debug().Msg("Make migrations")
	log.Debug().Msg("sql.open")
	log.Debug().Msg(dbURI)
	db, err := sql.Open("postgres", dbURI)
	log.Fatal().Err(err)

	log.Debug().Msg("set driver")
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	log.Fatal().Err(err)

	log.Debug().Msg("set instance")
	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/order/repository/postgres/migrations",
		"postgres",
		driver)
	log.Fatal().Err(err)

	log.Debug().Msg("m up")
	err = m.Up()
	log.Fatal().Err(err)

	log.Debug().Msg("migrate exit")
}
