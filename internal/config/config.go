package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Server Server
}

type Server struct {
	AppAddress      string        `env:"APP_PORT" env-default:"3000"`
	AppReadTimeout  time.Duration `env:"APP_READ_TIMEOUT" env-default:"10s"`
	AppWriteTimeout time.Duration `env:"APP_WRITE_TIMEOUT" env-default:"20s"`
	AppIdleTimeout  time.Duration `env:"APP_IDLE_TIMEOUT" env-default:"60s"`
}

// MustLoad returns Config in case no error
// If an error occurs, the app won't run and through a panic.
func MustLoad() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		panic("error loading environment variables")
	}

	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	cfg.Server.AppAddress = ":" + cfg.Server.AppAddress

	return &cfg
}
