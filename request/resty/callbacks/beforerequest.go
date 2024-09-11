package callbacks

import (
	"github.com/CrazyThursdayV50/pkgo/log"
	"github.com/CrazyThursdayV50/pkgo/trace"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

func TraceRequest(tracer trace.Tracer) func(*resty.Client, *resty.Request) error {
	return func(c *resty.Client, r *resty.Request) error {
		if tracer == nil {
			return nil
		}
		span, _ := tracer.NewSpan(r.Context())
		defer span.Finish()
		span.SetOperationName("outgoing request")
		span.SetTag("url", r.URL)
		// span.SetTag("path", r.RawRequest.URL.Path)
		span.SetTag("method", r.Method)
		span.SetTag("header", r.Header)
		span.SetTag("query", r.QueryParam.Encode())
		span.SetTag("body", r.Body)
		return nil
	}
}

func LogBeforeRequest(logger log.Logger) func(client *resty.Client, request *resty.Request) error {
	return func(c *resty.Client, r *resty.Request) error {
		logger.Info("outgoing request",
			zap.String("method", r.Method),
			// zap.String("path", r.RawRequest.URL.Path),
			zap.String("url", r.URL),
			zap.Any("header", r.Header),
			zap.String("query", r.QueryParam.Encode()),
			zap.Any("body", r.Body),
		)
		return nil
	}
}
