package rest

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/agatma/sprint1-http-server/internal/server/core/domain"
)

type MetricService interface {
	GetMetricValue(request *domain.MetricRequest) *domain.MetricResponse
	SetMetricValue(request *domain.SetMetricRequest) *domain.SetMetricResponse
	GetAllMetrics(request *domain.GetAllMetricsRequest) *domain.GetAllMetricsResponse
}

type handler struct {
	metricService MetricService
}

type API struct {
	srv *http.Server
}

func (a *API) Run() error {
	if err := a.srv.ListenAndServe(); err != nil {
		log.Printf("error occured during running server %v", err)
		return fmt.Errorf("failed run server: %w", err)
	}
	return nil
}

func NewAPI(metricService MetricService, cfg *Config) *API {
	h := &handler{
		metricService: metricService,
	}
	r := chi.NewRouter()
	r.Route("/update", func(r chi.Router) {
		r.Post("/{metricType}/{metricName}/{metricValue}", h.SetMetricValue)
	})
	r.Get("/value/{metricType}/{metricName}", h.GetMetricValue)
	r.Get("/", h.GetAllMetrics)
	return &API{
		srv: &http.Server{
			Addr:         cfg.Address,
			Handler:      r,
			ReadTimeout:  time.Second,
			WriteTimeout: time.Second,
		},
	}
}

func (h *handler) SetMetricValue(w http.ResponseWriter, req *http.Request) {
	metricType := chi.URLParam(req, "metricType")
	metricName := chi.URLParam(req, "metricName")
	metricValue := chi.URLParam(req, "metricValue")
	response := h.metricService.SetMetricValue(&domain.SetMetricRequest{
		MetricType:  metricType,
		MetricName:  metricName,
		MetricValue: metricValue,
	})
	if response.Error != nil {
		log.Printf(
			"failed to set metric value %s for metricType %s, metricName %s: %v",
			metricValue,
			metricType,
			metricName,
			response.Error,
		)
		switch {
		case errors.Is(response.Error, domain.ErrIncorrectMetricType):
			http.Error(w, response.Error.Error(), http.StatusBadRequest)
		case errors.Is(response.Error, domain.ErrIncorrectMetricValue):
			http.Error(w, response.Error.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "", http.StatusInternalServerError)
		}
	}
}

func (h *handler) GetMetricValue(w http.ResponseWriter, req *http.Request) {
	metricType, metricName := chi.URLParam(req, "metricType"), chi.URLParam(req, "metricName")
	response := h.metricService.GetMetricValue(&domain.MetricRequest{
		MetricType: metricType,
		MetricName: metricName,
	})
	if !response.Found {
		http.Error(w, domain.ErrItemNotFound.Error(), http.StatusNotFound)
		return
	}
	if response.Error != nil {
		log.Printf(
			"failed to get metric value for metricType %s, metricName %s: %v",
			metricType,
			metricName,
			response.Error,
		)
		if errors.Is(response.Error, domain.ErrIncorrectMetricType) {
			http.Error(w, domain.ErrItemNotFound.Error(), http.StatusNotFound)
		} else {
			http.Error(w, "", http.StatusInternalServerError)
		}
		return
	}
	if _, err := w.Write([]byte(response.MetricValue)); err != nil {
		return
	}
}

func (h *handler) GetAllMetrics(w http.ResponseWriter, req *http.Request) {
	gauge := h.metricService.GetAllMetrics(&domain.GetAllMetricsRequest{MetricType: domain.Gauge})
	if gauge.Error != nil {
		log.Printf("failed to get an item: %v for metricType %s", gauge.Error, domain.Gauge)
		http.Error(w, domain.ErrItemNotFound.Error(), http.StatusNotFound)
		return
	}
	counter := h.metricService.GetAllMetrics(&domain.GetAllMetricsRequest{MetricType: domain.Counter})
	if counter.Error != nil {
		log.Printf("failed to get an item: %v for metricType %s", gauge.Error, domain.Gauge)
		http.Error(w, domain.ErrItemNotFound.Error(), http.StatusNotFound)
		return
	}
	html := "<html><body><ul>"
	for key, value := range gauge.Values {
		html += fmt.Sprintf("<li>%s: %v</li>", key, value)
	}
	for key, value := range counter.Values {
		html += fmt.Sprintf("<li>%s: %v</li>", key, value)
	}
	html += "</ul></body></html>"
	w.Header().Set("Content-Type", "text/html")
	if _, err := w.Write([]byte(html)); err != nil {
		return
	}
}
