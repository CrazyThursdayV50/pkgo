package reconnector

import (
	"context"

	"github.com/CrazyThursdayV50/pkgo/log"
)

type ErrorCloserCheckerConnector[C ErrorCloserClosedChecker] func() (C, error)

func wrappedConnector[C ErrorCloserClosedChecker](f ErrorCloserCheckerConnector[C]) ConnectorFunc[C] {
	return func(ctx context.Context) (C, error) {
		return f()
	}
}

func NewWithErrorCloserCheckerConnector[C ErrorCloserClosedChecker](ctx context.Context, logger log.Logger, f func() (C, error)) *Reconnector[C] {
	return New(ctx, logger, wrappedConnector(f))
}
