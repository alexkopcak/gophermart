package postgres

import (
	"context"
	"errors"

	"github.com/alexkopcak/gophermart/internal/auth"
	"github.com/alexkopcak/gophermart/internal/models"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog/log"
)

type PostgresStorage struct {
	db *pgx.Conn
}

func NewPostgresStorage(dbURI string) auth.UserRepository {
	log.Debug().Msg("new postgres storage")
	MakeMigrations(dbURI)

	conn, err := pgx.Connect(context.Background(), dbURI)
	if err != nil {
		log.Fatal().Err(err)
	}
	return &PostgresStorage{
		db: conn,
	}
}

func (ps *PostgresStorage) CreateUser(ctx context.Context, user *models.User) error {
	log.Logger = log.With().Str("package", "postgres").Str("func", "CreateUser").Logger()

	log.Debug().Msg("enter")
	defer log.Debug().Msg("exit")

	log.Debug().Str("user", user.UserName).Msg("try to add user")

	_, err := ps.db.Exec(ctx,
		"INSERT INTO users "+
			"(login, password) "+
			"VALUES ($1, $2);", user.UserName, user.Password)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return auth.ErrUserAlreadyExsist
			}
		}
		log.Debug().Err(err).Msg("exit with error")
		return err
	}

	log.Debug().Msg("user created")
	return nil
}

func (ps *PostgresStorage) GetUser(ctx context.Context, userName string) (*models.User, error) {
	log.Logger = log.With().Str("package", "postgres").Str("func", "GetUser").Logger()

	log.Debug().Msg("enter")
	defer log.Debug().Msg("exit")

	log.Debug().Str("user", userName).Msg("get user by name")
	var user = new(models.User)
	err := ps.db.QueryRow(ctx,
		"SELECT login, password "+
			"FROM users "+
			"WHERE login = $1 "+
			"LIMIT 1;", userName).Scan(&user.UserName, &user.Password)
	if errors.Is(err, pgx.ErrNoRows) || user.UserName == "" {
		log.Debug().Str("user", userName).Msg("user not exsist")
		return nil, auth.ErrUserNotExsist
	}
	if err != nil {
		log.Err(err).Msg("exit with error")
		return nil, err
	}

	log.Debug().Str("user", userName).Msg("user finded at storage")
	return user, nil
}
