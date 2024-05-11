package workers

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v11"
)

const (
	defaultPollInterval   = 2
	defaultReportInterval = 10
)

type Config struct {
	Address        string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}

func NewConfig() (*Config, error) {
	var (
		cfg            Config
		flagRunAddr    *string
		pollInterval   *int
		reportInterval *int
	)
	flagRunAddr = flag.String("a", "localhost:8080", "run address")
	pollInterval = flag.Int("p", defaultPollInterval, " poll interval ")
	reportInterval = flag.Int("r", defaultReportInterval, " report interval ")
	flag.Parse()
	err := env.Parse(&cfg)
	if err != nil {
		return &cfg, fmt.Errorf("failed to get config for worker: %w", err)
	}
	if cfg.Address == "" {
		cfg.Address = *flagRunAddr
	}
	if cfg.ReportInterval == 0 {
		cfg.ReportInterval = *reportInterval
	}

	if cfg.PollInterval == 0 {
		cfg.PollInterval = *pollInterval
	}
	return &cfg, nil
}
