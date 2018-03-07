package model

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	RTotal     prometheus.Collector
	RUserIP    prometheus.Collector
	RRoute     prometheus.Collector
	RUserAgent prometheus.Collector
}

func CreateMetrics() *Metrics {
	return &Metrics{
		RTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "http_request_total",
			Help: "Requests total by Method, Status",
		}, []string{"method", "status"}),
		RUserIP: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "http_request_user_total",
			Help: "Requests total by Method, Status and User IP",
		}, []string{"method", "status", "ip"}),
		RRoute: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "http_request_route_total",
			Help: "Requests total by Method, Status and Route",
		}, []string{"method", "status", "route"}),
		RUserAgent: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "http_request_user_agent_total",
			Help: "Requests total by Method, Status, Route and User-Agent",
		}, []string{"method", "status", "route", "user_agent"}),
	}
}
