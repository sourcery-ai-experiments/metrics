package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
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

	log.Printf("made request %s. Got status code %d", resp.Request.URL, resp.StatusCode())
	return nil
}
