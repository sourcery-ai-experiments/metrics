package service

import (
	"strconv"

	"github.com/agatma/sprint1-http-server/internal/server/core/domain"
)

type MetricStorage interface {
	GetMetricValue(request *domain.MetricRequest) *domain.MetricResponse
	SetMetricValue(request *domain.SetMetricRequest) *domain.SetMetricResponse
	GetAllMetrics(request *domain.GetAllMetricsRequest) *domain.GetAllMetricsResponse
}

type MetricService struct {
	gaugeStorage   MetricStorage
	counterStorage MetricStorage
}

func NewMetricService(gauge MetricStorage, counter MetricStorage) *MetricService {
	return &MetricService{
		gaugeStorage:   gauge,
		counterStorage: counter,
	}
}

func (ms *MetricService) GetMetricValue(request *domain.MetricRequest) *domain.MetricResponse {
	switch request.MetricType {
	case domain.Gauge:
		return ms.gaugeStorage.GetMetricValue(request)
	case domain.Counter:
		return ms.counterStorage.GetMetricValue(request)
	default:
		return &domain.MetricResponse{
			Error: domain.ErrIncorrectMetricType,
		}
	}
}

func (ms *MetricService) SetMetricValue(request *domain.SetMetricRequest) *domain.SetMetricResponse {
	switch request.MetricType {
	case domain.Gauge:
		_, err := strconv.ParseFloat(request.MetricValue, 64)
		if err != nil {
			return &domain.SetMetricResponse{
				Error: domain.ErrIncorrectMetricValue,
			}
		}
		return ms.gaugeStorage.SetMetricValue(request)
	case domain.Counter:
		return ms.counterStorage.SetMetricValue(request)
	default:
		return &domain.SetMetricResponse{
			Error: domain.ErrIncorrectMetricType,
		}
	}
}

func (ms *MetricService) GetAllMetrics(request *domain.GetAllMetricsRequest) *domain.GetAllMetricsResponse {
	switch request.MetricType {
	case domain.Gauge:
		return ms.gaugeStorage.GetAllMetrics(request)
	case domain.Counter:
		return ms.counterStorage.GetAllMetrics(request)
	default:
		return &domain.GetAllMetricsResponse{
			Error: domain.ErrIncorrectMetricType,
		}
	}
}
