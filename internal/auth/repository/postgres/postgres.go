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

	log.Debug().Msg("pgx connect")
	conn, err := pgx.Connect(context.Background(), dbURI)
	if err != nil {
		log.Fatal().Err(err)
	}
	return &PostgresStorage{
		db: conn,
	}
}

func (ps *PostgresStorage) CreateUser(ctx context.Context, user *models.User) error {
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
		log.Debug().Str("error", err.Error()).Msg("got error")
		return err
	}

	return nil
}

func (ps *PostgresStorage) GetUser(ctx context.Context, userName string) (*models.User, error) {
	var user = new(models.User)
	err := ps.db.QueryRow(ctx, "SELECT login, password FROM users WHERE login = $1 ;", userName).Scan(&user.UserName, &user.Password)
	if errors.Is(err, pgx.ErrNoRows) || user.UserName == "" {
		return nil, auth.ErrUserNotExsist
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}
