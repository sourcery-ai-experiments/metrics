package memory

import (
	"strconv"
	"sync"

	"github.com/agatma/sprint1-http-server/internal/agent/core/domain"
)

type AgentMetricStorage struct {
	mux  *sync.Mutex
	data map[string]string
}

func NewAgentStorage(cfg *Config) *AgentMetricStorage {
	return &AgentMetricStorage{
		mux:  &sync.Mutex{},
		data: make(map[string]string),
	}
}

func (s *AgentMetricStorage) GetAllMetrics(req *domain.GetAllMetricsRequest) *domain.GetAllMetricsResponse {
	s.mux.Lock()
	defer s.mux.Unlock()
	return &domain.GetAllMetricsResponse{
		Values: s.data,
	}
}

func (s *AgentMetricStorage) GetMetricValue(req *domain.MetricRequest) *domain.MetricResponse {
	s.mux.Lock()
	defer s.mux.Unlock()
	value, found := s.data[req.MetricName]
	return &domain.MetricResponse{
		MetricValue: value,
		Found:       found,
	}
}

func (s *AgentMetricStorage) SetMetricValue(req *domain.SetMetricRequest) *domain.SetMetricResponse {
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

func setCounterMetricValue(req *domain.SetMetricRequest, s *AgentMetricStorage) *domain.SetMetricResponse {
	var currentValue int
	newValue, err := strconv.Atoi(req.MetricValue)
	if err != nil {
		return &domain.SetMetricResponse{
			Error: err,
		}
	}
	value, found := s.data[req.MetricName]
	if found {
		parsedValue, err := strconv.Atoi(value)
		if err != nil {
			return &domain.SetMetricResponse{
				Error: err,
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
