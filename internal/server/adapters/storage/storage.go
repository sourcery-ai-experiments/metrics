package storage

import (
	"errors"

	"github.com/agatma/sprint1-http-server/internal/server/adapters/storage/memory"
	"github.com/agatma/sprint1-http-server/internal/server/core/domain"
)

type MetricStorage interface {
	GetMetricValue(request *domain.MetricRequest) *domain.MetricResponse
	SetMetricValue(request *domain.SetMetricRequest) *domain.SetMetricResponse
	GetAllMetrics(request *domain.GetAllMetricsRequest) *domain.GetAllMetricsResponse
}

func NewStorage(conf Config) (MetricStorage, error) {
	if conf.Memory != nil {
		return memory.NewStorage(conf.Memory), nil
	}
	return nil, errors.New("no available storage")
}
