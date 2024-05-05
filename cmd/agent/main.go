package main

import (
	"fmt"
	"log"

	"github.com/agatma/sprint1-http-server/internal/agent/adapters/storage"
	"github.com/agatma/sprint1-http-server/internal/agent/adapters/storage/memory"
	"github.com/agatma/sprint1-http-server/internal/agent/adapters/workers"
	"github.com/agatma/sprint1-http-server/internal/agent/core/service"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg := workers.NewConfig()
	gaugeAgentStorage, err := storage.NewAgentStorage(storage.Config{
		Memory: &memory.Config{},
	})
	if err != nil {
		return fmt.Errorf("failed to initialize a storage: %w", err)
	}
	counterAgentStorage, err := storage.NewAgentStorage(storage.Config{
		Memory: &memory.Config{},
	})
	if err != nil {
		return fmt.Errorf("failed to initialize a storage: %w", err)
	}
	agentMetricService := service.NewAgentMetricService(gaugeAgentStorage, counterAgentStorage)
	worker := workers.NewAgentWorker(agentMetricService, cfg)
	if err := worker.Run(); err != nil {
		return fmt.Errorf("server has failed: %w", err)
	}
	return nil
}
