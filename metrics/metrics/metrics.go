package postgresql

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	ordersCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "orders_processed_total",
			Help: "Total number of processed orders.",
		},
	)

	// Добавьте еще кастомные метрики здесь
)

func init() {
	// Регистрация метрик в Prometheus
	prometheus.MustRegister(ordersCounter)
}

// Функция для увеличения счетчика выданных заказов
func IncrementOrdersCounter() {
	ordersCounter.Inc()
}
