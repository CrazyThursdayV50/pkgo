package gorm

import (
	"github.com/CrazyThursdayV50/pkgo/trace"
	"github.com/opentracing/opentracing-go/log"
	"gorm.io/gorm"
)

type (
	Handler     func(*gorm.DB)
	Interceptor func(op string, next Handler) Handler
)

func Sql(db *gorm.DB) string { return db.ToSQL(func(tx *gorm.DB) *gorm.DB { return tx }) }

func traceInterceptor(tracer trace.Tracer) func(string, Handler) Handler {
	return func(op string, next Handler) Handler {
		return func(db *gorm.DB) {
			if next != nil {
				defer next(db)
			}

			ctx := db.Statement.Context
			if ctx == nil {
				return
			}

			span, _ := tracer.NewSpanWithName(ctx, "gorm")
			defer span.Finish()

			var fields = []log.Field{
				log.Event(op),
				log.String("sql", Sql(db)),
				log.Int64("rows", db.RowsAffected),
			}

			if db.Error != nil {
				fields = append(fields, log.Error(db.Error))
			}

			span.LogFields(fields...)
		}
	}
}

func wrapInterceptors(interceptors ...Interceptor) Interceptor {
	var wrapped Interceptor
	for _, interceptor := range interceptors {
		if wrapped == nil {
			wrapped = interceptor
			continue
		}

		wrapped = func(op string, handler Handler) Handler {
			return interceptor(op, wrapped(op, handler))
		}
	}
	return wrapped
}

func registerInterceptors(db *gorm.DB, interceptors ...Interceptor) {
	var wrapped = wrapInterceptors(interceptors...)
	const InterceptorName = "GORM.Interceptor"

	{
		const name = "gorm:create"
		var callback = db.Callback().Create()
		callback.After(name).Register(InterceptorName, wrapped(name, nil))
	}

	{
		const name = "gorm:update"
		var callback = db.Callback().Update()
		callback.After(name).Register(InterceptorName, wrapped(name, nil))
	}

	{
		const name = "gorm:delete"
		var callback = db.Callback().Delete()
		callback.After(name).Register(InterceptorName, wrapped(name, nil))
	}

	{
		const name = "gorm:query"
		var callback = db.Callback().Query()
		callback.After(name).Register(InterceptorName, wrapped(name, nil))
	}

	{
		const name = "gorm:row"
		var callback = db.Callback().Row()
		callback.After(name).Register(InterceptorName, wrapped(name, nil))
	}

	{
		const name = "gorm:raw"
		var callback = db.Callback().Raw()
		callback.After(name).Register(InterceptorName, wrapped(name, nil))
	}
}
