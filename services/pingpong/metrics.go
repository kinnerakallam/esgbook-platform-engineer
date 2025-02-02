package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus"
)
var (
    pingRequestsReceived = prometheus.NewCounter(prometheus.CounterOpts{
        Name: "ping_requests_received_total",
        Help: "Total number of incoming /ping requests received",
    })

    pingFailuresTotal = prometheus.NewCounter(prometheus.CounterOpts{
        Name: "ping_failures_total",
        Help: "Total number of failed /ping requests",
    })
)

// initMetrics registers only metrics we want to track
func initMetrics() {
	prometheus.MustRegister(
        pingRequestsReceived,
        pingFailuresTotal,
    )
}

func startMetricsServer(cfg ConfigMetrics, wg *sync.WaitGroup) {
	defer wg.Done()

	port := fmt.Sprintf(":%v", cfg.Port)
	slog.With(slog.Any("port", port)).Info("metrics server started")
	http.Handle(cfg.Path, promhttp.Handler())
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Printf("error with listen and serve %v", err.Error())
	}
}
