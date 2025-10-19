package handlers

import (
	"github.com/Roflan4eg/auth-serivce/config"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricsHandler struct {
	handler http.Handler
}

func NewMetricsHandler(cfg *config.Config) *MetricsHandler {
	//return &MetricsHandler{
	//	handler: NewMetricsHandlerWithAuth(cfg),
	//}
	return &MetricsHandler{
		handler: promhttp.Handler(),
	}
}

func (m *MetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.handler.ServeHTTP(w, r)
}

func NewMetricsHandlerWithAuth(cfg *config.Config) *MetricsHandler {
	handler := promhttp.Handler()

	authenticatedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || user != cfg.Metrics.User || pass != cfg.Metrics.Password {
			w.Header().Set("WWW-Authenticate", `Basic realm="metrics"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		handler.ServeHTTP(w, r)
	})

	return &MetricsHandler{
		handler: authenticatedHandler,
	}
}
