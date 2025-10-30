package callbacks

import (
	"github.com/CrazyThursdayV50/pkgo/log"
	"github.com/CrazyThursdayV50/pkgo/trace"
	"github.com/go-resty/resty/v2"
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
		span.SetTag("method", r.Method)
		span.SetTag("header", r.Header)
		span.SetTag("query", r.QueryParam.Encode())
		span.SetTag("body", r.Body)
		return nil
	}
}

func LogBeforeRequest(logger log.Logger) func(client *resty.Client, request *resty.Request) error {
	return func(c *resty.Client, r *resty.Request) error {
		logger.Debugf("[outgoing request]Method(%s),URL(%s),Header(%v),Query(%s),Body(%v)",
			r.Method,
			r.URL,
			r.Header,
			r.QueryParam.Encode(),
			r.Body,
		)
		return nil
	}
}
