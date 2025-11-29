package main

import (
	"context"
	"github.com/gogf/gf/contrib/metric/otelmetric/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	"go.opentelemetry.io/otel/exporters/prometheus"
)

func main() {
	var ctx = gctx.New()
	// Set custom JSON logging handler to match google structured logging format.
	glog.SetDefaultHandler(LoggingJsonHandler)
	providerShutdown, _ := initMeterProvider()

	s := g.Server()
	s.BindHandler("/hello", func(r *ghttp.Request) {
		g.Log().Info(r.Context(), "hello world!!!")
		r.Response.Write("hello world")
	})
	// Prometheus metrics endpoint
	s.BindHandler("/metrics", otelmetric.PrometheusHandler)
	s.SetPort(8080)
	s.Run()

	defer providerShutdown(ctx)
}

func initMeterProvider() (func(context.Context) error, error) {
	// Set up OpenTelemetry Prometheus exporter to export metric as prometheus format.
	exporter, _ := prometheus.New()
	// OpenTelemetry provider.
	provider := otelmetric.MustProvider(
		otelmetric.WithReader(exporter),
	)
	provider.SetAsGlobal()
	return provider.Shutdown, nil
}
