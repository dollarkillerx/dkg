/**
*@program: vodka
*@description: prometheus 中间件
*@author: dollarkiller [dollarkiller@dollarkiller.com]
*@create: 2020-03-13 14:28
 */
package middleware

import (
	"context"
	prometheus2 "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc/status"
	"log"
	"net/http"
)

var Prometheus *prometheus

func init() {
	Prometheus = &prometheus{}
	Prometheus.init()
	log.Println("Prometheus Middleware Success")
}

type prometheus struct {
	requestCounter *prometheus2.CounterVec // 请求次数
	codeCounter    *prometheus2.CounterVec // 请求错误数量
	latencySummary *prometheus2.SummaryVec // 请求耗时分布
}

func (p *prometheus) init() {
	p.requestCounter = prometheus2.NewCounterVec(
		prometheus2.CounterOpts{
			Name: "vodka_server_request_total",
			Help: "Total number of RPCs completed on the server, regardless of success or failure.",
		},
		[]string{"service", "method"},
	)

	p.codeCounter = prometheus2.NewCounterVec(
		prometheus2.CounterOpts{
			Name: "vodka_server_handled_code_total",
			Help: "Total number of RPCs completed on the server, regardless of success or failure.",
		},
		[]string{"service", "method", "grpc_code"},
	)

	p.latencySummary = prometheus2.NewSummaryVec(
		prometheus2.SummaryOpts{
			Name:       "vodka_proc_cost",
			Help:       "RPC latency distributions.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"service", "method"},
	)

	prometheus2.MustRegister(p.requestCounter)
	prometheus2.MustRegister(p.codeCounter)
	prometheus2.MustRegister(p.latencySummary)
}

func (p *prometheus) Run(addr string) error {
	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(addr, nil)
}

func (p *prometheus) IncrRequest(ctx context.Context, serviceName, methodName string) {
	log.Println("aaa  111")
	p.requestCounter.WithLabelValues(serviceName, methodName).Inc()
}

func (p *prometheus) IncrCode(ctx context.Context, serviceName, methodName string, err error) {
	log.Println("aaa  122")

	st, _ := status.FromError(err)
	p.codeCounter.WithLabelValues(serviceName, methodName, st.Code().String()).Inc()
}

func (p *prometheus) Latency(ctx context.Context, serviceName, methodName string, us int64) {
	log.Println("aaa  1133")

	p.latencySummary.WithLabelValues(serviceName, methodName).Observe(float64(us))
}