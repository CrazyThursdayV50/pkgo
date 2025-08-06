package simple

import (
	"context"

	. "github.com/CrazyThursdayV50/pkgo/reconnector"
	. "github.com/CrazyThursdayV50/pkgo/reconnector/common"
)

type ReconnectorSimpleCloser[C Closer] Reconnector[*WrappedErrorCloserChecker[*WrappedErrorCloser[C]]]

func (r *ReconnectorSimpleCloser[Conn]) Reconnector() *Reconnector[*WrappedErrorCloserChecker[*WrappedErrorCloser[Conn]]] {
	return (*Reconnector[*WrappedErrorCloserChecker[*WrappedErrorCloser[Conn]]])(r)
}

func (r *ReconnectorSimpleCloser[Conn]) WithSimpleDialer(dialer func() Conn) *ReconnectorSimpleCloser[Conn] {
	var reconnector = r.Reconnector()
	reconnector = reconnector.WithSimpleDialer(func() *WrappedErrorCloserChecker[*WrappedErrorCloser[Conn]] {
		conn := dialer()
		return NewWrappedErrorCloserChecker(NewWrappedErrorCloser(conn))
	})
	return (*ReconnectorSimpleCloser[Conn])(reconnector)
}

func (r *ReconnectorSimpleCloser[Conn]) WithSimpleDialerContext(dialer func(context.Context) Conn) *ReconnectorSimpleCloser[Conn] {
	var reconnector = r.Reconnector()
	reconnector = reconnector.WithSimpleDialerContext(func(ctx context.Context) *WrappedErrorCloserChecker[*WrappedErrorCloser[Conn]] {
		conn := dialer(ctx)
		return NewWrappedErrorCloserChecker(NewWrappedErrorCloser(conn))
	})
	return (*ReconnectorSimpleCloser[Conn])(reconnector)
}

func (r *ReconnectorSimpleCloser[Conn]) WithDialer(dialer func() (Conn, error)) *ReconnectorSimpleCloser[Conn] {
	var reconnector = r.Reconnector()
	reconnector = reconnector.WithDialer(func() (*WrappedErrorCloserChecker[*WrappedErrorCloser[Conn]], error) {
		conn, err := dialer()
		return NewWrappedErrorCloserChecker(NewWrappedErrorCloser(conn)), err
	})
	return (*ReconnectorSimpleCloser[Conn])(reconnector)
}

func (r *ReconnectorSimpleCloser[Conn]) WithDialerContext(dialer func(context.Context) (Conn, error)) *ReconnectorSimpleCloser[Conn] {
	var reconnector = r.Reconnector()
	reconnector = reconnector.WithDialerContext(func(ctx context.Context) (*WrappedErrorCloserChecker[*WrappedErrorCloser[Conn]], error) {
		conn, err := dialer(ctx)
		return NewWrappedErrorCloserChecker(NewWrappedErrorCloser(conn)), err
	})
	return (*ReconnectorSimpleCloser[Conn])(reconnector)
}
