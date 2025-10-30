package callbacks

import (
	"github.com/CrazyThursdayV50/pkgo/log"
	"github.com/CrazyThursdayV50/pkgo/trace"
	"github.com/go-resty/resty/v2"
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
		logger.Debugf("[incoming response]Cost(%s),Status(%s),Method(%s),URL(%s),Result(%+v),Error(%v)",
			r.Time().String(),
			r.Status(),
			r.Request.Method,
			r.Request.URL,
			r.String(),
			r.Error(),
		)
		return nil
	}
}

func LogOnError(logger log.Logger) func(*resty.Request, error) {
	return func(r *resty.Request, err error) {
		logger.Errorf("[request failed]Method(%s),URL(%s),Header(%v),Query(%s),Body(%v),Result(%+v),Error(%v)",
			r.Method,
			r.URL,
			r.Header,
			r.QueryParam.Encode(),
			r.Body,
			r.Result,
			err,
		)
	}
}
