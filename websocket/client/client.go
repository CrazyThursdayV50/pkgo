package client

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/CrazyThursdayV50/pkgo/goo"
	"github.com/CrazyThursdayV50/pkgo/log"
	"github.com/CrazyThursdayV50/pkgo/reconnector"
	"github.com/CrazyThursdayV50/pkgo/reconnector/connection"
	"github.com/CrazyThursdayV50/pkgo/reconnector/dialerfunc"
	"github.com/CrazyThursdayV50/pkgo/websocket/compressor"

	"github.com/gorilla/websocket"
)

type wsreconnector = reconnector.Reconnector[*connection.WrappedChecker[*websocket.Conn]]

type (
	MessageHandler func(context.Context, log.Logger, int, []byte, func(error)) (int, []byte)
	Client         struct {
		ctx    context.Context
		cancel context.CancelFunc

		url          string
		l            log.Logger
		done         chan struct{}
		pingLoop     func(<-chan struct{}, *websocket.Conn)
		handler      MessageHandler
		pingHandler  func(string) error
		pongHandler  func(string) error
		onConnect    []func() (int, []byte)
		proxy        string
		readTimeout  time.Duration
		writeTimeout time.Duration

		c              compressor.Compressor
		enableCompress bool

		reconnector        *wsreconnector
		reconnectOnStartup bool
		doNotReconnect     bool
	}

	PingLoop func(done <-chan struct{}, conn *websocket.Conn)
)

const (
	TextMessage   = websocket.TextMessage
	BinaryMessage = websocket.BinaryMessage

	CloseMessage = websocket.CloseMessage
	PingMessage  = websocket.PingMessage
	PongMessage  = websocket.PongMessage
)

func (c *Client) listenClose() {
	goo.Go(func() {
		<-c.ctx.Done()
		c.l.Warn("context canceled")
		close(c.done)
		c.reconnector.Stop()
	})
}

func (c *Client) SendBinary(data []byte) error {
	return c.send(websocket.BinaryMessage, data)
}

func (c *Client) Send(data []byte) error {
	return c.send(websocket.TextMessage, data)
}

func (c *Client) Ping(data []byte) error {
	return c.send(websocket.PingMessage, data)
}

func (c *Client) Pong(data []byte) error {
	return c.send(websocket.PongMessage, data)
}

func (c *Client) send(typ int, data []byte) error {
	conn := c.reconnector.Connection()
	if c.writeTimeout != 0 {
		conn.Conn().SetWriteDeadline(time.Now().Add(c.writeTimeout))
	}

	switch typ {
	case websocket.CloseMessage:
		return conn.Conn().WriteControl(typ, data, time.Now().Add(time.Minute))

	case websocket.PingMessage:
		return conn.Conn().WriteControl(typ, data, time.Now().Add(time.Minute))

	case websocket.PongMessage:
		return conn.Conn().WriteControl(typ, data, time.Now().Add(time.Minute))

	case websocket.TextMessage:
		return conn.Conn().WriteMessage(typ, data)

	default:
		return conn.Conn().WriteMessage(typ, data)
	}
}

func (c *Client) readOnConn(conn *websocket.Conn, handler MessageHandler) error {
	if conn == nil {
		return nil
	}

	if c.readTimeout != 0 {
		conn.SetReadDeadline(time.Now().Add(c.readTimeout))
	}

	typ, data, err := conn.ReadMessage()
	if err != nil {
		return err
	}

	if c.c != nil {
		data, err = c.c.Uncompress(data)
		if err != nil {
			return err
		}
	}

	typ, message := handler(c.ctx, c.l, typ, data, func(err error) {
		c.l.Errorf("handle message failed. message: %v, error: %v", data, err)
	})

	switch typ {
	case websocket.BinaryMessage, websocket.TextMessage:
		if message == nil {
			return nil
		}

		if c.c != nil {
			message, err = c.c.Compress(message)
			if err != nil {
				return err
			}
		}

		return c.send(typ, message)

	default:

		c.send(typ, message)
	}

	return nil
}

func (c *Client) Run(ctx context.Context) error {
	c.reconnector.ReconnectOnStartup(c.reconnectOnStartup)
	c.reconnector.SetOnConnect(func(ctx context.Context, conn *connection.WrappedChecker[*websocket.Conn]) {
		if conn != nil {
			conn.Conn().SetPingHandler(c.pingHandler)
			conn.Conn().SetPongHandler(c.pongHandler)
			if c.pingLoop != nil {
				go c.pingLoop(c.ctx.Done(), conn.Conn())
			}

			c.l.Debugf("connect success to %s", c.url)
			for _, f := range c.onConnect {
				typ, data := f()
				if data != nil {
					err := c.send(typ, data)
					if err != nil {
						c.l.Errorf("send msg failed: %v", err)
						c.reconnector.Reconnect()
						return
					}
				}
			}

			go func() {
				conn := conn.Conn()
				if conn == nil {
					return
				}

				for {
					select {
					case <-c.done:
						c.l.Warn("EXIT")
						return

					case <-ctx.Done():
						c.l.Warn("EXIT")
						return

					default:
						err := c.readOnConn(conn, c.handler)
						if err != nil {
							c.l.Errorf("read msg failed: %v", err)
							c.reconnector.Reconnect()
							return
						}
					}
				}
			}()
		}
	})

	c.ctx, c.cancel = context.WithCancel(ctx)
	err := c.reconnector.Run(c.ctx)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Stop() {
	c.cancel()
}

func New(opts ...Option) *Client {
	var c Client
	c.done = make(chan struct{})

	for _, opt := range opts {
		opt(&c)
	}

	dialer := websocket.DefaultDialer
	dialer.EnableCompression = c.enableCompress
	switch c.proxy {
	case "":
	case "env":
		dialer.Proxy = http.ProxyFromEnvironment
	default:
		url, err := url.Parse(c.proxy)
		if err == nil {
			dialer.Proxy = http.ProxyURL(url)
		}
	}

	dialerFunc := dialerfunc.CloserDialerContext[*websocket.Conn](func(ctx context.Context) (*websocket.Conn, error) {
		c.l.Debugf("connect to %s ...", c.url)
		conn, _, err := dialer.DialContext(ctx, c.url, nil)
		if err != nil {
			return nil, err
		}

		return conn, nil
	})

	c.reconnector = reconnector.New(dialerFunc.Wrap()).WithLogger(c.l)
	if c.doNotReconnect {
		c.reconnector.DoNotReconnect()
	}

	c.reconnector.ReconnectInterval(time.Second)
	c.listenClose()

	return &c
}

func (c *Client) UpdateOptions(opts ...Option) {
	for _, opt := range opts {
		opt(c)
	}
}
