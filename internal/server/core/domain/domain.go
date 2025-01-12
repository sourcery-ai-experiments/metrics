package domain

import "errors"

const (
	Gauge   = "gauge"
	Counter = "counter"
)

var (
	ErrIncorrectMetricType  = errors.New("incorrect metric value")
	ErrIncorrectMetricValue = errors.New("incorrect metric value")
	ErrItemNotFound         = errors.New("item not found")
)

type MetricRequest struct {
	MetricType string
	MetricName string
}

type MetricResponse struct {
	MetricValue string
	Found       bool
	Error       error
}

type SetMetricRequest struct {
	MetricType  string
	MetricName  string
	MetricValue string
}

type SetMetricResponse struct {
	Error error
}

type GetAllMetricsRequest struct {
	MetricType string
}

type GetAllMetricsResponse struct {
	Values map[string]string
	Error  error
}
