package rest

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Address string `env:"ADDRESS"`
}

func NewConfig() *Config {
	var flagRunAddr *string
	var cfg Config
	flagRunAddr = flag.String("a", ":8080", "address and port to run server")
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	flag.Parse()
	if cfg.Address == "" {
		cfg.Address = *flagRunAddr
	}
	return &cfg
}
