package connection

import "io"

type Checker interface {
	io.Closer
	IsClosed() bool
}

type WrappedChecker[Conn io.Closer] struct {
	closed bool
	conn   Conn
}

func (c *WrappedChecker[Conn]) Conn() Conn { return c.conn }
func (c *WrappedChecker[Conn]) IsClosed() bool {
	return c.closed
}

func (c *WrappedChecker[Conn]) Close() error {
	if c.closed {
		return nil
	}

	err := c.Conn().Close()
	c.closed = err == nil
	return err
}

func WrapToChecker[Conn io.Closer](conn Conn) *WrappedChecker[Conn] {
	return &WrappedChecker[Conn]{conn: conn}
}
