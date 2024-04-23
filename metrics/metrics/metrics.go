package postgresql

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	OrdersCounter       prometheus.Counter
	OrdersInProgress    prometheus.Gauge
	ProcessingHistogram prometheus.Histogram
)

func Initialize() {
	OrdersCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "orders_count",
		Help: "Total number of orders processed",
	})

	OrdersInProgress = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "orders_in_progress",
		Help: "Number of orders currently in progress",
	})

	ProcessingHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "order_processing_duration",
		Help:    "Order processing duration in seconds",
		Buckets: prometheus.LinearBuckets(0, 10, 5), // 5 интервалов по 10 секунд
	})

	prometheus.MustRegister(OrdersCounter, OrdersInProgress, ProcessingHistogram)
}
