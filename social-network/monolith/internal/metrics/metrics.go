package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	loadTestWrites *prometheus.CounterVec
}

func NewMetrics() *Metrics {
	m := &Metrics{
		loadTestWrites: prometheus.NewCounterVec(
			prometheus.CounterOpts{ //nolint:exhaustruct // остальное по-умолчанию.
				Namespace: "load_test",
				Name:      "writes_count_total",
				Help:      "Counter of load test writes",
			},
			[]string{
				"success",
			},
		),
	}

	prometheus.MustRegister(m.loadTestWrites)

	return m
}

func (m *Metrics) IncLoadTestWrites(success bool) {
	labels := prometheus.Labels{"success": "true"}

	if !success {
		labels["success"] = "false"
	}

	m.loadTestWrites.With(labels).Inc()
}
