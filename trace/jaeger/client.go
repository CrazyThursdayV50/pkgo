package jaeger

import (
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
)

var DefaultLogger = jaeger.StdLogger

func InitJaeger(cfg *Config, logger jaeger.Logger) (opentracing.Tracer, io.Closer, error) {
	jaegerCfgInstance := jaegercfg.Configuration{
		ServiceName: cfg.ServiceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           cfg.LogSpans,
			LocalAgentHostPort: cfg.Host,
			CollectorEndpoint:  "http://" + cfg.Host + "/api/traces",
		},
	}

	return jaegerCfgInstance.NewTracer(
		jaegercfg.Logger(logger),
		jaegercfg.Metrics(metrics.NullFactory),
	)
}
