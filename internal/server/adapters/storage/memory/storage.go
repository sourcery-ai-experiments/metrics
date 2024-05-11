package memory

import (
	"strconv"
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
	if req.MetricType == domain.Counter {
		return setCounterMetricValue(req, s)
	}
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

func setCounterMetricValue(req *domain.SetMetricRequest, s *MetricStorage) *domain.SetMetricResponse {
	var currentValue int
	newValue, err := strconv.Atoi(req.MetricValue)
	if err != nil {
		return &domain.SetMetricResponse{
			Error: domain.ErrIncorrectMetricValue,
		}
	}
	value, found := s.data[req.MetricName]
	if found {
		parsedValue, err := strconv.Atoi(value)
		if err != nil {
			return &domain.SetMetricResponse{
				Error: domain.ErrIncorrectMetricValue,
			}
		}
		currentValue = parsedValue
	}
	newValue += currentValue
	s.data[req.MetricName] = strconv.Itoa(newValue)
	return &domain.SetMetricResponse{
		Error: nil,
	}
}
