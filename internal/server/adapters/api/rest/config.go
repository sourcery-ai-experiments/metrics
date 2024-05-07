package rest

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Address string `env:"ADDRESS"`
}

func NewConfig() (*Config, error) {
	var flagRunAddr *string
	var cfg Config
	flagRunAddr = flag.String("a", ":8080", "address and port to run server")
	err := env.Parse(&cfg)
	if err != nil {
		return &cfg, fmt.Errorf("failed to get config for server: %w", err)
	}
	flag.Parse()
	if cfg.Address == "" {
		cfg.Address = *flagRunAddr
	}
	return &cfg, nil
}
