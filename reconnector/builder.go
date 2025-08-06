package reconnector

/*
import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/CrazyThursdayV50/pkgo/log"
)

// reconnWrapper is an internal adapter that wraps the user's connection object (of any type C)
// and makes it conform to the ErrorCloserClosedChecker interface, which the core Reconnector logic uses.
type reconnWrapper[C any] struct {
	conn         C
	isClosedFunc func(C) bool
}

// Close intelligently calls the close method on the wrapped connection.
// It uses type assertions to check for `Close() error` or `Close()` signatures.
func (w *reconnWrapper[C]) Close() error {
	switch c := any(w.conn).(type) {
	case interface{ Close() error }:
		return c.Close()
	case io.Closer: // Double-check for io.Closer
		return c.Close()
	case interface{ Close() }:
		c.Close()
		return nil
	default:
		// If no Close method is found, do nothing.
		// The builder logic ensures a valid closer is provided.
		return nil
	}
}

// Closed checks if the connection is closed.
// It prioritizes a user-provided checker function, then looks for standard `Closed()` or `IsClosed()` methods.
// If none are found, it defaults to false.
func (w *reconnWrapper[C]) Closed() bool {
	if w.isClosedFunc != nil {
		return w.isClosedFunc(w.conn)
	}

	type closedChecker interface{ Closed() bool }
	type isClosedChecker interface{ IsClosed() bool }

	switch c := any(w.conn).(type) {
	case closedChecker:
		return c.Closed()
	case isClosedChecker:
		return c.IsClosed()
	default:
		return false
	}
}

// Config holds the configuration for the Reconnector, built using Options.
type Config[C any] struct {
	logger            log.Logger
	reconnectInterval time.Duration
	connector         func(context.Context) (C, error)
	isClosedFunc      func(C) bool
}

// Option defines a function that configures a Reconnector.
type Option[C any] func(*Config[C])

// WithLogger provides a logger for the Reconnector.
func WithLogger[C any](logger log.Logger) Option[C] {
	return func(c *Config[C]) {
		c.logger = logger
	}
}

// WithReconnectInterval sets the time to wait between reconnection attempts.
func WithReconnectInterval[C any](interval time.Duration) Option[C] {
	return func(c *Config[C]) {
		c.reconnectInterval = interval
	}
}

// WithConnector sets the connection factory function (without context).
func WithConnector[C any](f func() (C, error)) Option[C] {
	return func(c *Config[C]) {
		c.connector = func(_ context.Context) (C, error) {
			return f()
		}
	}
}

// WithConnectorContext sets the connection factory function (with context).
func WithConnectorContext[C any](f func(context.Context) (C, error)) Option[C] {
	return func(c *Config[C]) {
		c.connector = f
	}
}

// WithSimpleConnector sets a connection factory that doesn't return an error.
func WithSimpleConnector[C any](f func() C) Option[C] {
	return func(c *Config[C]) {
		c.connector = func(_ context.Context) (C, error) {
			return f(), nil
		}
	}
}

// WithSimpleConnectorContext sets a connection factory with context that doesn't return an error.
func WithSimpleConnectorContext[C any](f func(context.Context) C) Option[C] {
	return func(c *Config[C]) {
		c.connector = func(ctx context.Context) (C, error) {
			return f(ctx), nil
		}
	}
}

// WithIsClosedChecker provides a custom function to check if the connection is closed.
// This is useful when the connection object doesn't have a standard `IsClosed() bool` or `Closed() bool` method.
func WithIsClosedChecker[C any](f func(C) bool) Option[C] {
	return func(c *Config[C]) {
		c.isClosedFunc = f
	}
}

// NewReconnector creates and initializes a new Reconnector with a fluent, option-based configuration.
// It is highly flexible and can adapt to various connection functions and connection object types.
//
// C is the generic type of the user's connection object.
func NewReconnector[C any](ctx context.Context, opts ...Option[C]) (*Reconnector[*reconnWrapper[C]], error) {
	// 1. Setup default config
	cfg := &Config[C]{
		logger:            log.NewNop(),
		reconnectInterval: 3 * time.Second,
	}

	// 2. Apply all user-provided options
	for _, opt := range opts {
		opt(cfg)
	}

	// 3. Validate essential options
	if cfg.connector == nil {
		return nil, fmt.Errorf("reconnector: connector function is required. Use WithConnector or a similar option")
	}

	// 4. Create the normalized connector function that the core Reconnector will use.
	// This function calls the user's connector and wraps the result in our internal adapter.
	normalizedConnector := func(ctx context.Context) (*reconnWrapper[C], error) {
		rawConn, err := cfg.connector(ctx)
		if err != nil {
			return nil, err
		}

		// Ensure the provided connection type has a valid Close method.
		// This is a runtime check to fail fast if the connection is not closable.
		if _, ok := any(rawConn).(interface{ Close() error }); !ok {
			if _, ok := any(rawConn).(interface{ Close() }); !ok {
				return nil, fmt.Errorf("reconnector: connection type %T must have a 'Close() error' or 'Close()' method", rawConn)
			}
		}

		return &reconnWrapper[C]{
			conn:         rawConn,
			isClosedFunc: cfg.isClosedFunc,
		}, nil
	}

	// 5. Create the core Reconnector instance
	r := New(ctx, cfg.logger, normalizedConnector)
	r.ReconnectInterval(cfg.reconnectInterval)

	return r, nil
}

*/
