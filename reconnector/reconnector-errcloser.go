package reconnector

import (
	"context"

	"github.com/CrazyThursdayV50/pkgo/log"
)

type ErrorCloserConnector[C ErrorCloser] func() (C, error)

func wrappedErrorCloserCheckerConnector[C ErrorCloser](f ErrorCloserConnector[C]) ErrorCloserCheckerConnector[*WrappedErrorCloserChecker[C]] {
	return func() (*WrappedErrorCloserChecker[C], error) {
		closer, err := f()
		errcloser := NewWrappedErrorCloserChecker(closer)
		return errcloser, err
	}
}
func NewWithErrorCloserConnector[C ErrorCloser](ctx context.Context, logger log.Logger, f func() (C, error)) *Reconnector[*WrappedErrorCloserChecker[C]] {
	return NewWithErrorCloserCheckerConnector(ctx, logger, wrappedErrorCloserCheckerConnector(f))
}

type ErrorCloserConnectorContext[C ErrorCloser] func(context.Context) (C, error)

func wrappedErrorCloserCheckerConnectorContext[C ErrorCloser](f ErrorCloserConnectorContext[C]) ConnectorFunc[*WrappedErrorCloserChecker[C]] {
	return func(ctx context.Context) (*WrappedErrorCloserChecker[C], error) {
		closer, err := f(ctx)
		errcloser := NewWrappedErrorCloserChecker(closer)
		return errcloser, err
	}
}

func NewWithErrorCloserConnectorContext[C ErrorCloser](ctx context.Context, logger log.Logger, f func(context.Context) (C, error)) *Reconnector[*WrappedErrorCloserChecker[C]] {
	return New(ctx, logger, wrappedErrorCloserCheckerConnectorContext(f))
}
