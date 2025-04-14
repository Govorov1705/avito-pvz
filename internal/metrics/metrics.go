package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	HttpRequestsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Общее количество HTTP-запросов",
	})

	HttpRequestDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Продолжительность HTTP-запросов",
		Buckets: prometheus.DefBuckets,
	})

	PvzCreatedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "pvz_created_total",
		Help: "Количество созданных ПВЗ",
	})

	ReceptionsCreatedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "receptions_created_total",
		Help: "Количество созданных приемок",
	})

	ProductsAddedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "products_added_total",
		Help: "Количество добавленных товаров",
	})
)

func StartMetricsServer() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":9000", nil)
	}()
}
