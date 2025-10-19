package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	GRPCRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "grpc_request_duration_seconds",
		Help:    "Duration of gRPC requests in seconds",
		Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1, 5},
	}, []string{"method", "status_code"})

	GRPCRequestCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "grpc_requests_total",
		Help: "Total number of gRPC requests",
	}, []string{"method", "status_code"})

	GRPCRequestErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "grpc_request_errors_total",
		Help: "Total number of gRPC errors by type",
	}, []string{"method", "error_type"})

	//GRPCRequestSize = promauto.NewHistogramVec(prometheus.HistogramOpts{
	//	Name:    "grpc_request_size_bytes",
	//	Help:    "Size of gRPC requests",
	//	Buckets: prometheus.ExponentialBuckets(100, 10, 5),
	//}, []string{"method"})

	GRPCRequestsInFlight = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "grpc_requests_in_flight",
		Help: "Current number of gRPC requests being processed",
	}, []string{"service", "method"})
)
