package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Manager struct {
	executionsCount      *prometheus.CounterVec
	executionsRunning    *prometheus.GaugeVec
	executionsWaiting    *prometheus.GaugeVec
	executionsPending    *prometheus.GaugeVec
	httpRequestsTotal    *prometheus.CounterVec
	httpRequestDuration  *prometheus.HistogramVec
	httpRequestsInFlight *prometheus.GaugeVec
}

func NewManager() *Manager {
	return &Manager{
		executionsCount: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "flowctl",
			Name:      "executions_total",
			Help:      "Total processed executions",
		},
			[]string{"namespace", "flow_id", "state"},
		),
		executionsRunning: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "flowctl",
			Name:      "executions_running",
			Help:      "Number of running transactions",
		},
			[]string{"namespace", "flow_id"},
		),
		executionsWaiting: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "flowctl",
			Name:      "executions_waiting",
			Help:      "Number of executions waiting for user action",
		},
			[]string{"namespace", "flow_id"},
		),
		executionsPending: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "flowctl",
			Name:      "executions_pending",
			Help:      "Number of executions in pending state",
		},
			[]string{"namespace", "flow_id"},
		),
		httpRequestsTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "flowctl",
			Name:      "http_requests_total",
			Help:      "Total number of HTTP requests",
		},
			[]string{"method", "path", "status"},
		),
		httpRequestDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "flowctl",
			Name:      "http_request_duration_seconds",
			Help:      "HTTP request duration in seconds",
			Buckets:   prometheus.DefBuckets,
		},
			[]string{"method", "path", "status"},
		),
		httpRequestsInFlight: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "flowctl",
			Name:      "http_requests_in_flight",
			Help:      "Number of HTTP requests currently being processed",
		},
			[]string{"method", "path"},
		),
	}
}

func (m *Manager) Register() {
	prometheus.MustRegister(
		m.executionsCount,
		m.executionsRunning,
		m.executionsWaiting,
		m.executionsPending,
		m.httpRequestsTotal,
		m.httpRequestDuration,
		m.httpRequestsInFlight,
	)
}

func (m *Manager) GetHandler() http.Handler {
	return promhttp.Handler()
}

func (m *Manager) IncrementExecutionCount(namespace, flowID, state string) {
	m.executionsCount.WithLabelValues(namespace, flowID, state).Inc()
}

func (m *Manager) SetExecutionsRunning(namespace, flowID string, value float64) {
	m.executionsRunning.WithLabelValues(namespace, flowID).Set(value)
}

func (m *Manager) IncExecutionsRunning(namespace, flowID string) {
	m.executionsRunning.WithLabelValues(namespace, flowID).Inc()
}

func (m *Manager) DecExecutionsRunning(namespace, flowID string) {
	m.executionsRunning.WithLabelValues(namespace, flowID).Dec()
}

func (m *Manager) SetExecutionsWaiting(namespace, flowID string, value float64) {
	m.executionsWaiting.WithLabelValues(namespace, flowID).Set(value)
}

func (m *Manager) IncExecutionsWaiting(namespace, flowID string) {
	m.executionsWaiting.WithLabelValues(namespace, flowID).Inc()
}

func (m *Manager) DecExecutionsWaiting(namespace, flowID string) {
	m.executionsWaiting.WithLabelValues(namespace, flowID).Dec()
}

func (m *Manager) SetExecutionsPending(namespace, flowID string, value float64) {
	m.executionsPending.WithLabelValues(namespace, flowID).Set(value)
}

func (m *Manager) HTTPMetricsMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			req := c.Request()
			path := c.Path()
			method := req.Method

			m.httpRequestsInFlight.WithLabelValues(method, path).Inc()
			defer m.httpRequestsInFlight.WithLabelValues(method, path).Dec()

			err := next(c)

			duration := time.Since(start).Seconds()
			status := strconv.Itoa(c.Response().Status)

			m.httpRequestsTotal.WithLabelValues(method, path, status).Inc()
			m.httpRequestDuration.WithLabelValues(method, path, status).Observe(duration)

			return err
		}
	}
}
