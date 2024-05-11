package domain

const (
	Gauge       = "gauge"
	Counter     = "counter"
	PollCount   = "PollCount"
	RandomValue = "RandomValue"
)

type Metrics struct {
	Values map[string]string
}

type MetricRequest struct {
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
