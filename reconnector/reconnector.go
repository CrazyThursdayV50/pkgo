package reconnector

import (
	"context"
	"time"

	"github.com/CrazyThursdayV50/pkgo/log"
	"github.com/CrazyThursdayV50/pkgo/reconnector/connection"
)

type ConnectorFunc[Conn connection.Checker] func(ctx context.Context) (Conn, error)
type Reconnector[Conn connection.Checker] struct {
	ctx    context.Context
	cancel context.CancelFunc
	logger log.Logger

	reconnectOnStartup  bool
	reconnectInterval   time.Duration
	reconnectSignalChan chan struct{}
	sendReconnectSignal func()
	newConn             ConnectorFunc[Conn]
	conn                Conn
	onConnect           func(context.Context, Conn)
}

func (r *Reconnector[Conn]) WithLogger(logger log.Logger) *Reconnector[Conn] {
	r.logger = logger
	return r
}

// func (r *Reconnector[Conn]) WithContext(ctx context.Context) *Reconnector[Conn] {
// r.ctx, r.cancel = context.WithCancel(ctx)
// return r
// }

func New[Conn connection.Checker](dialer func(context.Context) (Conn, error)) *Reconnector[Conn] {
	var r Reconnector[Conn]
	r.newConn = dialer
	r.reconnectSignalChan = make(chan struct{}, 1)
	r.sendReconnectSignal = func() {
		select {
		case <-r.reconnectSignalChan:
		default:
			r.reconnectSignalChan <- struct{}{}
		}
	}
	return &r
}

func (r *Reconnector[Conn]) Stop() {
	r.cancel()
}

func (r *Reconnector[Conn]) connect() error {
	ctx := context.WithoutCancel(r.ctx)
	conn, err := r.newConn(ctx)
	if err != nil {
		return err
	}

	r.conn = conn
	if r.onConnect != nil {
		r.onConnect(ctx, conn)
	}

	return nil
}

func (r *Reconnector[Conn]) Run(ctx context.Context) error {
	r.ctx, r.cancel = context.WithCancel(ctx)

	go func() {
		for {
			select {
			case <-r.ctx.Done():
				conn := r.conn
				if !conn.IsClosed() {
					err := conn.Close()
					if err != nil {
						r.logger.Errorf("reconnector closed. close connection failed: %v", err)
						return
					}
				}

				r.logger.Warn("exit reconnector and connection CLOSED")
				return

			case <-r.reconnectSignalChan:
				conn := r.conn
				if !conn.IsClosed() {
					err := conn.Close()
					if err != nil {
						r.logger.Errorf("Close connection failed: %v", err)
					}
				}

				err := r.connect()
				if err != nil {
					r.logger.Errorf("Connect failed: %v", err)
					r.logger.Debugf("Reconnect in %s", r.reconnectInterval.String())
					time.Sleep(r.reconnectInterval)
					r.sendReconnectSignal()
				}
			}
		}
	}()

	switch r.reconnectOnStartup {
	case false:
		err := r.connect()
		if err != nil {
			return err
		}

	default:
		r.sendReconnectSignal()
	}

	return nil
}

func (r *Reconnector[Conn]) Connection() Conn {
	return r.conn
}

func (r *Reconnector[Conn]) SetOnConnect(onConnect func(context.Context, Conn)) {
	r.onConnect = onConnect
}

func (r *Reconnector[Conn]) ReconnectOnStartup(ok bool) {
	r.reconnectOnStartup = ok
}

func (r *Reconnector[Conn]) ReconnectInterval(interval time.Duration) {
	r.reconnectInterval = interval
}

func (r *Reconnector[Conn]) Reconnect() {
	r.sendReconnectSignal()
}
