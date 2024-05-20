package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"

	"github.com/agatma/sprint1-http-server/internal/agent/logger"
)

func SendMetrics(host string, metricType string, metricName string, metricValue string) error {
	client := resty.New()
	resp, err := client.R().
		SetRawPathParams(map[string]string{
			"metricType":  metricType,
			"metricName":  strings.ToLower(metricName),
			"metricValue": metricValue,
		}).
		Post(host + "/update/{metricType}/{metricName}/{metricValue}")

	if err != nil {
		return fmt.Errorf("failed to send metrics: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("bad request. Status Code %d", resp.StatusCode())
	}

	logger.Log.Info(
		"made http request",
		zap.String("uri", resp.Request.URL),
		zap.String("method", resp.Request.Method),
		zap.Int("statusCode", resp.StatusCode()),
		zap.Duration("duration", resp.Time()),
	)
	return nil
}
