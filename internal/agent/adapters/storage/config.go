package storage

import "github.com/agatma/sprint1-http-server/internal/agent/adapters/storage/memory"

type Config struct {
	Memory *memory.Config
}
