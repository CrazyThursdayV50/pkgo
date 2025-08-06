package reconnector

import (
	"context"
	"time"

	"github.com/CrazyThursdayV50/pkgo/log"
)

type ReconnectorInterface[Conn ErrorCloserClosedChecker] interface {
	WithContext(context.Context) ReconnectorInterface[Conn]
	WithLogger(log.Logger) ReconnectorInterface[Conn]
	WithSimpleDialer(dialer func() Conn) ReconnectorInterface[Conn]
	WithSimpleDialerContext(dialer func(context.Context) Conn) ReconnectorInterface[Conn]
	WithDialer(dialer func() (Conn, error)) ReconnectorInterface[Conn]
	WithDialerContext(dialer func(ctx context.Context) (Conn, error)) ReconnectorInterface[Conn]
	Stop()
	Run() error
	Connection() Conn
	SetOnConnect(func(context.Context, Conn))
	ReconnectInterval(time.Duration)
	Reconnect()
}

// type ReconnectorCloser[C io.Closer] Reconnector[*WrappedErrorCloserChecker[C]]
