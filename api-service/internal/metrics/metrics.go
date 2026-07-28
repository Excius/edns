package metrics

import "github.com/prometheus/client_golang/prometheus"

type HTTPMetrics struct {
	RequestsTotal      prometheus.CounterVec
	RequestDuration    prometheus.HistogramVec
	RequestsInProgress prometheus.Gauge
}

type BusinessMetrics struct {
	UsersCreated         prometheus.Counter
	NotificationsCreated prometheus.Counter
}

type Metrics struct {
	HTTP     HTTPMetrics
	Business BusinessMetrics
}

func NewMetrics() *Metrics {

	m := &Metrics{
		HTTP: HTTPMetrics{
			RequestsTotal: *prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Name: "api_http_requests_total",
					Help: "Total HTTP requests.",
				},
				[]string{"method", "route", "status"},
			),
			RequestDuration: *prometheus.NewHistogramVec(
				prometheus.HistogramOpts{
					Name:    "api_http_request_duration_seconds",
					Help:    "HTTP request duration.",
					Buckets: prometheus.DefBuckets,
				},
				[]string{"method", "route", "status"},
			),
			RequestsInProgress: prometheus.NewGauge(
				prometheus.GaugeOpts{
					Name: "api_http_requests_in_progress",
					Help: "Current HTTP requests in progress.",
				},
			),
		},

		Business: BusinessMetrics{
			UsersCreated: prometheus.NewCounter(
				prometheus.CounterOpts{
					Name: "api_users_created_total",
					Help: "Total users created.",
				},
			),
			NotificationsCreated: prometheus.NewCounter(
				prometheus.CounterOpts{
					Name: "api_notifications_created_total",
					Help: "Total notifications successfully created.",
				},
			),
		},
	}

	prometheus.MustRegister(

		m.HTTP.RequestsTotal,
		m.HTTP.RequestDuration,
		m.HTTP.RequestsInProgress,

		m.Business.UsersCreated,
		m.Business.NotificationsCreated,
	)

	return m
}
