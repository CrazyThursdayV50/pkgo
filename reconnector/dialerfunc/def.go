package dialerfunc

import (
	"context"
	"io"

	"github.com/CrazyThursdayV50/pkgo/reconnector/connection"
)

type CheckerBuilder[Conn connection.Checker] func() Conn
type CheckerBuilderContext[Conn connection.Checker] func(context.Context) Conn
type CheckerDialer[Conn connection.Checker] func() (Conn, error)
type CheckerDialerContext[Conn connection.Checker] func(context.Context) (Conn, error)

type CloserBuilder[Conn io.Closer] func() Conn
type CloserBuilderContext[Conn io.Closer] func(context.Context) Conn
type CloserDialer[Conn io.Closer] func() (Conn, error)
type CloserDialerContext[Conn io.Closer] func(context.Context) (Conn, error)

type SimpleCloserBuilder[Conn connection.SimpleCloser] func() Conn
type SimpleCloserBuilderContext[Conn connection.SimpleCloser] func(context.Context) Conn
type SimpleCloserDialer[Conn connection.SimpleCloser] func() (Conn, error)
type SimpleCloserDialerContext[Conn connection.SimpleCloser] func(context.Context) (Conn, error)
