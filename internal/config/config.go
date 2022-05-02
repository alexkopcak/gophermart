package config

import (
	"flag"

	"github.com/caarlos0/env"
)

type Config struct {
	RunAddress           string `env:"RUN_ADDRESS" envDefault:"127.0.0.1:8080"`
	DataBaseURI          string `env:"DATABASE_URI" envDefault:"postgres://user:pass@localhost:5432/gophermart?sslmode=disable"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS" envDefault:"127.0.0.1:8081"`
	HashSalt             string `env:"HASH_SALT" envDefault:"hash salt"`
	SigningKey           string `env:"SIGNING_KEY" envDefault:"signing key"`
	TokenTTL             int    `env:"TOKEN_TTL" envDefault:"600"`
}

func Init() *Config {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}

	cfg.RunAddress = *(flag.String("a", cfg.RunAddress, "Server address"))
	cfg.DataBaseURI = *(flag.String("d", cfg.DataBaseURI, "Database URI"))
	cfg.AccrualSystemAddress = *(flag.String("r", cfg.AccrualSystemAddress, "Accrual system address"))

	return &cfg
}
