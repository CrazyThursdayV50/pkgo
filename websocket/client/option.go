package client

import (
	"time"

	"github.com/CrazyThursdayV50/pkgo/log"
	"github.com/CrazyThursdayV50/pkgo/websocket/compressor"
)

type Option func(*Client)

func WithProxy(proxy string) Option {
	return func(c *Client) { c.proxy = proxy }
}
func WithMessageHandler(handler MessageHandler) Option {
	return func(c *Client) { c.handler = handler }
}

func WithURL(url string) Option {
	return func(c *Client) { c.url = url }
}

func WithLogger(logger log.Logger) Option {
	return func(c *Client) { c.l = logger }
}

func WithPingLoop(f PingLoop) Option {
	return func(c *Client) { c.pingLoop = f }
}

func WithSendOnConnect(f func() (int, []byte)) Option {
	return func(c *Client) { c.onConnect = append(c.onConnect, f) }
}

func WithCompressor(compressor compressor.Compressor) Option {
	return func(c *Client) { c.c = compressor }
}

func WithDefaultCompress(ok bool) Option {
	return func(c *Client) { c.enableCompress = ok }
}

func WithPingHandler(timeout time.Duration, f func(string) error) Option {
	return func(c *Client) {
		c.pingHandler = func(appData string) error {
			if f == nil {
				return c.reconnector.Connection().Conn().WriteControl(PongMessage, []byte(appData), time.Now().Add(timeout))
			}

			return f(appData)
		}
	}
}

func WithPongHandler(timeout time.Duration, f func(string) error) Option {
	return func(c *Client) {
		c.pongHandler = func(appData string) error {
			if f == nil {
				return nil
			}

			return f(appData)
		}
	}
}

func WithReconnectOnStartup(reconnectOnStartup bool) Option {
	return func(c *Client) { c.reconnectOnStartup = reconnectOnStartup }
}
