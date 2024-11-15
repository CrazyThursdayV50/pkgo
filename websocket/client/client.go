package client

import (
	"context"
	"time"

	"github.com/CrazyThursdayV50/gotils/pkg/async/goo"
	"github.com/CrazyThursdayV50/pkgo/log"
	"github.com/CrazyThursdayV50/pkgo/websocket/client/compressor"
	"go.uber.org/zap"

	"github.com/gorilla/websocket"
)

type (
	MessageHandler func(context.Context, log.Logger, []byte, func(error)) []byte
	Client         struct {
		ctx       context.Context
		cancel    context.CancelFunc
		url       string
		conn      *websocket.Conn
		l         log.Logger
		done      chan struct{}
		pingLoop  func(<-chan struct{}, *websocket.Conn)
		handler   MessageHandler
		onConnect []func() []byte

		c compressor.Compressor
	}
)

func (c *Client) listenClose() {
	goo.Go(func() {
		<-c.ctx.Done()
		close(c.done)
		c.conn.Close()
	})
}

func (c *Client) connect() error {
	dialer := websocket.DefaultDialer
	conn, _, err := dialer.DialContext(c.ctx, c.url, nil)
	if err != nil {
		return err
	}

	c.conn = conn
	c.l.Info("connect success")

	for _, f := range c.onConnect {
		data := f()
		if data != nil {
			err = c.send(data)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Client) reconnect() error {
	time.Sleep(time.Second)
	return c.connect()
}

func (c *Client) start() {
	var connected bool
	for !connected {
		err := c.connect()
		if err == nil {
			connected = true
			continue
		}

		c.l.Error("connect failed, try to reconnect", zap.Error(err))
		timer := time.NewTimer(time.Second * 5)
		select {
		case <-c.done:
			return

		case <-timer.C:
			_ = c.connect()
		}
	}
}

func (c *Client) Send(data []byte) error {
	return c.send(data)
}

func (c *Client) send(data []byte) error {
	c.l.Debug("send", zap.ByteString("message", data))
	return c.conn.WriteMessage(websocket.TextMessage, data)
}

func (c *Client) read(handler MessageHandler) error {
	_, data, err := c.conn.ReadMessage()
	if err != nil {
		return err
	}

	message := handler(c.ctx, c.l, data, func(err error) {
		c.l.Error("handle message failed", zap.Any("message", data), zap.Error(err))
	})

	if message != nil {
		return c.send(message)
	}

	return nil
}

func (c *Client) runPingLoop() {
	if c.pingLoop == nil {
		return
	}

	goo.Go(func() {
		c.pingLoop(c.done, c.conn)
	})
}

func (c *Client) Run() {
	c.start()
	c.runPingLoop()

	goo.Go(func() {
		for {
			select {
			case <-c.done:
				c.l.Info("exit")
				return

			default:
				err := c.read(c.handler)
				if err != nil {
					c.l.Error("read message failed", zap.Error(err))
					_ = c.reconnect()
				}
			}
		}
	})
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

	c.listenClose()
	return &c
}
