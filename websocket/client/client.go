package client

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/CrazyThursdayV50/pkgo/goo"
	"github.com/CrazyThursdayV50/pkgo/log"
	"github.com/CrazyThursdayV50/pkgo/websocket/compressor"
	"go.uber.org/zap"

	"github.com/gorilla/websocket"
)

type (
	MessageHandler func(context.Context, log.Logger, int, []byte, func(error)) (int, []byte)
	Client         struct {
		ctx         context.Context
		cancel      context.CancelFunc
		url         string
		conn        *websocket.Conn
		l           log.Logger
		done        chan struct{}
		pingLoop    func(<-chan struct{}, *websocket.Conn)
		handler     MessageHandler
		pingHandler func(string) error
		pongHandler func(string) error
		onConnect   []func() (int, []byte)
		proxy       string

		c                   compressor.Compressor
		enableCompress      bool
		writeControlTimeout time.Duration
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
		c.conn.Close()
	})
}

func (c *Client) connect() error {
	dialer := websocket.DefaultDialer
	dialer.EnableCompression = c.enableCompress
	if c.proxy != "" {
		url, err := url.Parse(c.proxy)
		if err == nil {
			dialer.Proxy = http.ProxyURL(url)
		}
	}

	conn, _, err := dialer.DialContext(c.ctx, c.url, nil)
	if err != nil {
		return err
	}

	conn.SetPingHandler(c.pingHandler)
	conn.SetPongHandler(c.pongHandler)

	c.conn = conn
	c.l.Info("connect success")

	for _, f := range c.onConnect {
		typ, data := f()
		if data != nil {
			err = c.send(typ, data)
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

// func (c *Client) start() {
// 	var connected bool
// 	for !connected {
// 		err := c.connect()
// 		if err == nil {
// 			connected = true
// 			continue
// 		}

// 		c.l.Errorf("connect failed, try to reconnect: %v", zap.Error(err))
// 		timer := time.NewTimer(time.Second * 5)
// 		select {
// 		case <-c.done:
// 			return

//			case <-timer.C:
//				_ = c.connect()
//			}
//		}
//	}

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
	switch typ {
	case websocket.CloseMessage:
		c.l.Debugf("send CLOSE")
		return c.conn.WriteControl(typ, data, time.Now().Add(time.Minute))

	case websocket.PingMessage:
		c.l.Debugf("send PING")
		return c.conn.WriteControl(typ, data, time.Now().Add(time.Minute))

	case websocket.PongMessage:
		c.l.Debugf("send PONG")
		return c.conn.WriteControl(typ, data, time.Now().Add(time.Minute))

	case websocket.TextMessage:
		c.l.Debugf("send: %v", zap.String("message", string(data)))
		return c.conn.WriteMessage(typ, data)

	default:
		c.l.Debugf("send: %v", zap.ByteString("message", data))
		return c.conn.WriteMessage(typ, data)
	}
}

func (c *Client) read(handler MessageHandler) error {
	typ, data, err := c.conn.ReadMessage()
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
		c.l.Error("handle message failed", zap.Any("message", data), zap.Error(err))
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

func (c *Client) runPingLoop() {
	if c.pingLoop == nil {
		return
	}

	goo.Go(func() {
		c.pingLoop(c.done, c.conn)
	})
}

func (c *Client) onMessage() {
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

func (c *Client) Run() error {
	err := c.connect()
	if err != nil {
		return err
	}

	c.runPingLoop()
	c.onMessage()
	return nil
}

func (c *Client) Stop() {
	c.cancel()
}

func New(opts ...Option) *Client {
	var c Client
	c.done = make(chan struct{})
	c.writeControlTimeout = time.Second

	for _, opt := range opts {
		opt(&c)
	}

	c.listenClose()
	return &c
}

func (c *Client) UpdateOptions(opts ...Option) {
	for _, opt := range opts {
		opt(c)
	}
}
