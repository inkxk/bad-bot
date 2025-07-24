package config

import (
	"log"
	"time"

	env "github.com/caarlos0/env/v11"
)

func NewConfig() (*Config, error) {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		log.Print(err)
		return nil, err
	}

	return &cfg, nil
}

type Config struct {
	LogLevel        string        `env:"LOG_LEVEL" envDefault:"info"`
	GracefulTimeout time.Duration `env:"GRACEFUL_TIMEOUT" envDefault:"30s"`
	HTTPServer      HTTPServer
	Line            Line
}

type HTTPServer struct {
	Port              string        `env:"HTTP_SERVER_PORT"`
	ReadHeaderTimeout time.Duration `env:"HTTP_SERVER_READ_HEADER_TIMEOUT" envDefault:"10s"`
}

type Line struct {
	LineChannelSecret     string `env:"LINE_CHANNEL_SECRET"`
	LineChanneAccessToken string `env:"LINE_CHANNEL_ACCESS_TOKEN"`
}
