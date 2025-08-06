package connection

import "io"

type SimpleCloser interface {
	Close()
}

var _ io.Closer = (*WrappedCloser[SimpleCloser])(nil)

type WrappedCloser[Conn SimpleCloser] struct {
	conn Conn
}

func (c *WrappedCloser[Conn]) Conn() Conn { return c.conn }
func (c *WrappedCloser[Conn]) Close() error {
	c.conn.Close()
	return nil
}

func WrapToCloser[Conn SimpleCloser](conn Conn) *WrappedCloser[Conn] {
	return &WrappedCloser[Conn]{conn: conn}
}
