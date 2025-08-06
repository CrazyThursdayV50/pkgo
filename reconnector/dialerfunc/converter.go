package dialerfunc

import (
	"context"
	"io"

	"github.com/CrazyThursdayV50/pkgo/reconnector/connection"
)

func builderToDialerContext[Conn any](builder func() Conn) func(context.Context) (Conn, error) {
	return func(context.Context) (Conn, error) {
		return builder(), nil
	}
}

func builderContextToDialerContext[Conn any](builder func(context.Context) Conn) func(context.Context) (Conn, error) {
	return func(ctx context.Context) (Conn, error) {
		return builder(ctx), nil
	}
}

func dialerToContext[Conn any](dialer func() (Conn, error)) func(context.Context) (Conn, error) {
	return func(ctx context.Context) (Conn, error) {
		return dialer()
	}
}

func closerToCheckerContext[Conn io.Closer](dialer func(context.Context) (Conn, error)) func(context.Context) (*connection.WrappedChecker[Conn], error) {
	return func(ctx context.Context) (*connection.WrappedChecker[Conn], error) {
		conn, err := dialer(ctx)
		return connection.WrapToChecker(conn), err
	}
}

func simpleCloserToCheckerContext[Conn connection.SimpleCloser](dialer func(context.Context) (Conn, error)) func(context.Context) (*connection.WrappedChecker[*connection.WrappedCloser[Conn]], error) {
	return closerToCheckerContext(func(ctx context.Context) (*connection.WrappedCloser[Conn], error) {
		conn, err := dialer(ctx)
		return connection.WrapToCloser(conn), err
	})
}

func (c CheckerDialer[Conn]) Wrap() CheckerDialerContext[Conn] {
	return dialerToContext(c)
}

func (c CheckerBuilderContext[Conn]) Wrap() CheckerDialerContext[Conn] {
	return builderContextToDialerContext(c)
}

func (c CheckerBuilder[Conn]) Wrap() CheckerDialerContext[Conn] {
	return builderToDialerContext(c)
}

func (c CloserDialerContext[Conn]) Wrap() CheckerDialerContext[*connection.WrappedChecker[Conn]] {
	return closerToCheckerContext(c)
}

func (c CloserDialer[Conn]) Wrap() CheckerDialerContext[*connection.WrappedChecker[Conn]] {
	return closerToCheckerContext(dialerToContext(c))
}

func (c CloserBuilderContext[Conn]) Wrap() CheckerDialerContext[*connection.WrappedChecker[Conn]] {
	return closerToCheckerContext(builderContextToDialerContext(c))
}

func (c CloserBuilder[Conn]) Wrap() CheckerDialerContext[*connection.WrappedChecker[Conn]] {
	return closerToCheckerContext(builderToDialerContext(c))
}

func (c SimpleCloserDialerContext[Conn]) Wrap() CheckerDialerContext[*connection.WrappedChecker[*connection.WrappedCloser[Conn]]] {
	return simpleCloserToCheckerContext(c)
}

func (c SimpleCloserDialer[Conn]) Wrap() CheckerDialerContext[*connection.WrappedChecker[*connection.WrappedCloser[Conn]]] {
	return simpleCloserToCheckerContext(dialerToContext(c))
}

func (c SimpleCloserBuilderContext[Conn]) Wrap() CheckerDialerContext[*connection.WrappedChecker[*connection.WrappedCloser[Conn]]] {
	return simpleCloserToCheckerContext(builderContextToDialerContext(c))
}

func (c SimpleCloserBuilder[Conn]) Wrap() CheckerDialerContext[*connection.WrappedChecker[*connection.WrappedCloser[Conn]]] {
	return simpleCloserToCheckerContext(builderToDialerContext(c))
}
