package storage

import (
	"errors"

	"github.com/agatma/sprint1-http-server/internal/agent/adapters/storage/memory"
	"github.com/agatma/sprint1-http-server/internal/agent/core/domain"
)

type AgentMetricStorage interface {
	GetMetricValue(request *domain.MetricRequest) *domain.MetricResponse
	SetMetricValue(request *domain.SetMetricRequest) *domain.SetMetricResponse
	GetAllMetrics(request *domain.GetAllMetricsRequest) *domain.GetAllMetricsResponse
}

func NewAgentStorage(conf Config) (AgentMetricStorage, error) {
	if conf.Memory != nil {
		return memory.NewAgentStorage(conf.Memory), nil
	}
	return nil, errors.New("no available agent storage")
}
