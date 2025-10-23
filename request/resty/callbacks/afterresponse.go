package callbacks

import (
	"github.com/CrazyThursdayV50/pkgo/log"
	"github.com/CrazyThursdayV50/pkgo/trace"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

func TraceResponse(tracer trace.Tracer) func(*resty.Client, *resty.Response) error {
	return func(c *resty.Client, r *resty.Response) error {
		if tracer == nil {
			return nil
		}
		span, _ := tracer.NewSpan(r.Request.Context())
		defer span.Finish()
		span.SetOperationName("incoming response")
		span.SetTag("url", r.Request.URL)
		span.SetTag("method", r.Request.Method)
		span.SetTag("status", r.Status())
		span.SetTag("cost", r.Time().String())
		span.SetTag("result", r.Result())
		span.SetTag("error", r.Error())
		return nil
	}
}

func LogAfterResponse(logger log.Logger) func(*resty.Client, *resty.Response) error {
	return func(c *resty.Client, r *resty.Response) error {
		logger.Info("incoming response",
			zap.String("method", r.Request.Method),
			zap.String("url", r.Request.URL),
			zap.String("status", r.Status()),
			zap.String("cost", r.Time().String()),
			zap.Any("result", r.Result()),
			zap.Any("error", r.Error()),
		)
		return nil
	}
}

func LogOnError(logger log.Logger) func(*resty.Request, error) {
	return func(r *resty.Request, err error) {
		logger.Error("request failed",
			zap.String("method", r.Method),
			zap.String("url", r.URL),
			zap.Any("header", r.Header),
			zap.String("query", r.QueryParam.Encode()),
			zap.Any("body", r.Body),
			zap.Any("result", r.Result),
			zap.Error(err),
		)
	}
}
