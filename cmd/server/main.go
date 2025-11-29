package main

import (
	"context"
	"github.com/gogf/gf/contrib/metric/otelmetric/v2"
	"github.com/gogf/gf/contrib/trace/otlphttp/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gmetric"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
)

const (
	instrument        = "github.com/this/is/a/otel/go/demo"
	instrumentVersion = "v1.0"
	serviceName       = "otlp-http-server"          // Name of the OTel service for tracing
	endpoint          = "host.docker.internal:4318" // OTel tracing endpoint
	path              = "/v1/traces"                // OTel tracing path
)

func main() {
	var ctx = gctx.New()
	var shutdown, _ = otlphttp.Init(serviceName, endpoint, path)
	// Set custom JSON logging handler to match google structured logging format.
	glog.SetDefaultHandler(LoggingJsonHandler)
	providerShutdown, _ := initMeterProvider()
	counter, gauge, histogram := initMetrics()

	s := g.Server()
	s.BindHandler("/hello", func(r *ghttp.Request) {
		g.Log().Info(r.Context(), "hello world!!!")
		addMetricValue(r.Context(), counter, gauge, histogram)
		addManualTrace(r.Context())
		r.Response.Write("hello world")
	})
	// Prometheus metrics endpoint
	s.BindHandler("/metrics", otelmetric.PrometheusHandler)
	s.SetPort(8080)
	s.Run()

	defer providerShutdown(ctx)
	defer shutdown(ctx)
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

func initMetrics() (gmetric.Counter, gmetric.UpDownCounter, gmetric.Histogram) {
	meter := gmetric.GetGlobalProvider().Meter(gmetric.MeterOption{
		Instrument:        instrument,
		InstrumentVersion: instrumentVersion,
	})
	counter := meter.MustCounter(
		"goframe.metric.demo.counter",
		gmetric.MetricOption{
			Help: "This is a simple demo for Counter usage",
			Unit: "bytes",
		},
	)
	gauge := meter.MustUpDownCounter(
		"goframe.metric.demo.gauge",
		gmetric.MetricOption{
			Help: "This is a simple demo for UpDownCounter usage",
			Unit: "%",
		},
	)
	histogram := meter.MustHistogram(
		"goframe.metric.demo.histogram",
		gmetric.MetricOption{
			Help:    "This is a simple demo for histogram usage",
			Unit:    "ms",
			Buckets: []float64{0, 10, 20, 50, 100, 500, 1000, 2000, 5000, 10000},
		},
	)
	return counter, gauge, histogram
}

func addMetricValue(ctx context.Context, counter gmetric.Counter, gauge gmetric.UpDownCounter, histogram gmetric.Histogram) {
	counter.Add(ctx, 1)

	gauge.Add(ctx, 10) // Add adds the given value to the counter. It panics if the value is < 0
	gauge.Dec(ctx)

	histogram.Record(1)
	histogram.Record(20)
	histogram.Record(30)
	histogram.Record(101)
	histogram.Record(2000)
	histogram.Record(9000)
	histogram.Record(20000)
}

func addManualTrace(ctx context.Context) {
	tracer := otel.GetTracerProvider().Tracer("otel-go-demo")
	ctx, span := tracer.Start(ctx, "hello-outer-span")
	span.AddEvent("hello-event")
	traceInner(ctx)
	defer span.End()
}

func traceInner(ctx context.Context) {
	_, span := gtrace.NewSpan(ctx, "hello-inner-span")
	defer span.End()
}
