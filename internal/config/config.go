package config

import (
	"flag"

	"github.com/caarlos0/env"
)

type Config struct {
	RunAddress           string `env:"RUN_ADDRESS" envDefault:"127.0.0.1:8080"`
	DataBaseURI          string `env:"DATABASE_URI" envDefault:"postgresql://localhost:5432/gophermart?sslmode=disable"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS" envDefault:"127.0.0.1:8081"`
	SecretKey            string `env:"SECRET_KEY" envDefault:"my Secret Key"`
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
