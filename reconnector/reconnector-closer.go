package reconnector

import (
	"context"

	"github.com/CrazyThursdayV50/pkgo/log"
)

type CloserConnector[C Closer] func() (C, error)

func wrappedErrorCloserConnector[C Closer](f CloserConnector[C]) ErrorCloserConnector[*WrappedErrorCloser[C]] {
	return func() (*WrappedErrorCloser[C], error) {
		closer, err := f()
		errcloser := NewWrappedErrorCloser(closer)
		return errcloser, err
	}
}

func NewWithCloserConnector[C Closer](
	ctx context.Context,
	logger log.Logger,
	f func() (C, error),
) *Reconnector[*WrappedErrorCloserChecker[*WrappedErrorCloser[C]]] {
	c := NewWithErrorCloserConnector(ctx, logger, wrappedErrorCloserConnector(f))
	return c
}

type CloserConnectorContext[C Closer] func(context.Context) (C, error)

func wrappedErrorCloserConnectorContext[C Closer](f CloserConnectorContext[C]) ErrorCloserConnectorContext[*WrappedErrorCloser[C]] {
	return func(ctx context.Context) (*WrappedErrorCloser[C], error) {
		closer, err := f(ctx)
		errcloser := NewWrappedErrorCloser(closer)
		return errcloser, err
	}
}

func NewWithCloserConnectorContext[C Closer](
	ctx context.Context,
	logger log.Logger,
	f func(context.Context) (C, error),
) *Reconnector[*WrappedErrorCloserChecker[*WrappedErrorCloser[C]]] {
	c := NewWithErrorCloserConnectorContext(ctx, logger, wrappedErrorCloserConnectorContext(f))
	return c
}
