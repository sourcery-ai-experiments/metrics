package rest

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/agatma/sprint1-http-server/internal/server/adapters/api/middleware"
	"github.com/agatma/sprint1-http-server/internal/server/core/domain"
	"github.com/agatma/sprint1-http-server/internal/server/logger"
)

const (
	metricType  = "metricType"
	metricValue = "metricValue"
	metricName  = "metricName"
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
		logger.Log.Error("error occured during running server: ", zap.Error(err))
		return fmt.Errorf("failed run server: %w", err)
	}
	return nil
}

func NewAPI(metricService MetricService, cfg *Config) *API {
	h := &handler{
		metricService: metricService,
	}
	r := chi.NewRouter()
	r.Use(middleware.RequestLogging)
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
	mType := chi.URLParam(req, metricType)
	mName := chi.URLParam(req, metricName)
	mValue := chi.URLParam(req, metricValue)
	response := h.metricService.SetMetricValue(&domain.SetMetricRequest{
		MetricType:  mType,
		MetricName:  mName,
		MetricValue: mValue,
	})
	if response.Error != nil {
		logger.Log.Error("failed to set metric",
			zap.String(metricValue, mValue),
			zap.String(metricType, mType),
			zap.String(metricName, mName),
			zap.Error(response.Error),
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
	w.WriteHeader(http.StatusOK)
}

func (h *handler) GetMetricValue(w http.ResponseWriter, req *http.Request) {
	mType, mName := chi.URLParam(req, metricType), chi.URLParam(req, metricName)
	response := h.metricService.GetMetricValue(&domain.MetricRequest{
		MetricType: mType,
		MetricName: mName,
	})
	if !response.Found {
		http.Error(w, domain.ErrItemNotFound.Error(), http.StatusNotFound)
		return
	}
	if response.Error != nil {
		logger.Log.Error("failed to get metric",
			zap.String(metricType, mType),
			zap.String(metricName, mName),
			zap.Error(response.Error),
		)
		if errors.Is(response.Error, domain.ErrIncorrectMetricType) {
			http.Error(w, domain.ErrItemNotFound.Error(), http.StatusNotFound)
		} else {
			http.Error(w, "", http.StatusInternalServerError)
		}
		return
	}
	if _, err := w.Write([]byte(response.MetricValue)); err != nil {
		w.WriteHeader(http.StatusOK)
		return
	}
}

func (h *handler) GetAllMetrics(w http.ResponseWriter, req *http.Request) {
	gauge := h.metricService.GetAllMetrics(&domain.GetAllMetricsRequest{MetricType: domain.Gauge})
	if gauge.Error != nil {
		logger.Log.Error(
			"failed to get an item",
			zap.String(metricType, domain.Gauge),
			zap.Error(gauge.Error),
		)
		http.Error(w, domain.ErrItemNotFound.Error(), http.StatusNotFound)
		return
	}
	counter := h.metricService.GetAllMetrics(&domain.GetAllMetricsRequest{MetricType: domain.Counter})
	if counter.Error != nil {
		logger.Log.Error(
			"failed to get an item",
			zap.String(metricType, domain.Counter),
			zap.Error(counter.Error),
		)
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
		w.WriteHeader(http.StatusOK)
		return
	}
}
