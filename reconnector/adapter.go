package reconnector

/*

import "io"

// Adapter is a helper struct to adapt different connection types to the
// ErrorCloserClosedChecker interface required by the Reconnector.
// It provides a default implementation for Closed() which always returns false.
// You can embed this struct into your custom connection wrapper to easily
// satisfy the interface.
//
// Example:
//
//	type MyConnection struct {
//	    // ... your connection fields
//	}
//
//	func (c *MyConnection) Close() error {
//	    // ... close logic
//	}
//
//	// Define an adapter for MyConnection
//	type MyConnectionAdapter struct {
//	    *MyConnection
//	    Adapter
//	}
//
//	// Now, MyConnectionAdapter implicitly satisfies ErrorCloserClosedChecker
//
//	func connect() (*MyConnectionAdapter, error) {
//	    conn, err := makeNewMyConnection()
//	    if err != nil {
//	        return nil, err
//	    }
//	    return &MyConnectionAdapter{MyConnection: conn}, nil
//	}
//
//	// Use it with the reconnector
//	reconn := reconnector.New(ctx, logger, func(ctx context.Context) (*MyConnectionAdapter, error) {
//	    return connect()
//	})
type Adapter struct{}

// Closed always returns false.
// For connections that have a state, you should override this method
// in your own adapter implementation to reflect the actual connection state.
func (a Adapter) Closed() bool {
	return false
}

// CloserAdapter wraps an io.Closer (which has a `Close() error` method) and
// adapts it to the ErrorCloserClosedChecker interface.
type CloserAdapter[T io.Closer] struct {
	Adapter
	Conn T
}

// NewCloserAdapter creates a new adapter for any type that implements io.Closer.
func NewCloserAdapter[T io.Closer](conn T) *CloserAdapter[T] {
	return &CloserAdapter[T]{Conn: conn}
}

// Close calls the underlying connection's Close method.
func (a *CloserAdapter[T]) Close() error {
	return a.Conn.Close()
}

// SimpleCloserAdapter wraps a type that has just a `Close()` method (no error return)
// and adapts it to the ErrorCloserClosedChecker interface.
type SimpleCloserAdapter[T Closer] struct {
	Adapter
	Conn T
}

// NewSimpleCloserAdapter creates a new adapter for any type that implements the
// simple Closer interface (with just a Close() method).
func NewSimpleCloserAdapter[T Closer](conn T) *SimpleCloserAdapter[T] {
	return &SimpleCloserAdapter[T]{Conn: conn}
}

// Close calls the underlying connection's Close method and always returns nil.
func (a *SimpleCloserAdapter[T]) Close() error {
	a.Conn.Close()
	return nil
}

*/
