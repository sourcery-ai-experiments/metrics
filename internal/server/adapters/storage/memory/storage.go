package memory

import (
	"sync"

	"github.com/agatma/sprint1-http-server/internal/server/core/domain"
)

type MetricStorage struct {
	mux  *sync.Mutex
	data map[string]string
}

func NewStorage(cfg *Config) *MetricStorage {
	return &MetricStorage{
		mux:  &sync.Mutex{},
		data: make(map[string]string),
	}
}

func (s *MetricStorage) GetMetricValue(req *domain.MetricRequest) *domain.MetricResponse {
	s.mux.Lock()
	defer s.mux.Unlock()
	value, found := s.data[req.MetricName]
	return &domain.MetricResponse{
		MetricValue: value,
		Found:       found,
	}
}

func (s *MetricStorage) SetMetricValue(req *domain.SetMetricRequest) *domain.SetMetricResponse {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.data[req.MetricName] = req.MetricValue
	return &domain.SetMetricResponse{
		Error: nil,
	}
}

func (s *MetricStorage) GetAllMetrics(req *domain.GetAllMetricsRequest) *domain.GetAllMetricsResponse {
	s.mux.Lock()
	defer s.mux.Unlock()
	return &domain.GetAllMetricsResponse{
		Values: s.data,
	}
}
