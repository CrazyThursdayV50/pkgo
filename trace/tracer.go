package trace

import (
	"context"

	"github.com/opentracing/opentracing-go"
)

type Tracer interface {
	NewSpan(ctx context.Context) (opentracing.Span, context.Context)
	NewSpanWithName(ctx context.Context, name string) (opentracing.Span, context.Context)
}
