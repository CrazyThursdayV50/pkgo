package jaeger

import (
	"context"
	"fmt"
	"io"
	"runtime"

	"github.com/CrazyThursdayV50/gotils/pkg/async/goo"
	"github.com/CrazyThursdayV50/pkgo/log"
	jaegerlogger "github.com/CrazyThursdayV50/pkgo/log/jaeger"
	"github.com/CrazyThursdayV50/pkgo/trace"
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

type Tracer struct {
	t opentracing.Tracer
}

type tracer struct {
	name string
	t    opentracing.Tracer
}

func (t *Tracer) NewTracer(name string) trace.Tracer {
	return &tracer{name: name, t: t.t}
}

func (t *tracer) NewSpan(ctx context.Context) (opentracing.Span, context.Context) {
	pc, _, _, _ := runtime.Caller(1)
	callerName := runtime.FuncForPC(pc).Name()
	return opentracing.StartSpanFromContextWithTracer(ctx, t.t, fmt.Sprintf("%s.%s", t.name, callerName))
}

func (t *tracer) NewSpanWithName(ctx context.Context, name string) (opentracing.Span, context.Context) {
	return opentracing.StartSpanFromContextWithTracer(ctx, t.t, fmt.Sprintf("%s.%s", t.name, name))
}

func New(ctx context.Context, cfg *Config, logger log.Logger) (*Tracer, error) {
	jl := jaegerlogger.New(logger)
	tracer, closer, err := InitJaeger(cfg, jl)
	if err != nil {
		return nil, err
	}

	goo.Go(func() {
		<-ctx.Done()
		closer.Close()
	})

	return &Tracer{t: tracer}, nil
}
