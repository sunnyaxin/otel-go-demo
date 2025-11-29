package main

import (
	"github.com/gogf/gf/contrib/metric/otelmetric/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"go.opentelemetry.io/otel/exporters/prometheus"
)

func main() {
	var ctx = gctx.New()
	// Prometheus exporter to export metrics as Prometheus format.
	exporter, _ := prometheus.New()
	// OpenTelemetry provider.
	provider := otelmetric.MustProvider(
		otelmetric.WithReader(exporter),
	)
	provider.SetAsGlobal()
	defer provider.Shutdown(ctx)

	s := g.Server()
	s.BindHandler("/hello", func(r *ghttp.Request) {
		r.Response.Write("hello world")
	})
	s.BindHandler("/metrics", otelmetric.PrometheusHandler)
	s.SetPort(8080)
	s.Run()
}
