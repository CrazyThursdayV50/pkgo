package client

import (
	"context"

	"github.com/CrazyThursdayV50/pkgo/log"
	"github.com/CrazyThursdayV50/pkgo/websocket/compressor"
	"github.com/gorilla/websocket"
)

type Option func(*Client)

func WithMessageHandler(handler MessageHandler) Option {
	return func(c *Client) { c.handler = handler }
}

func WithURL(url string) Option {
	return func(c *Client) { c.url = url }
}

func WithLogger(logger log.Logger) Option {
	return func(c *Client) { c.l = logger }
}

func WithPingLoop(f func(done <-chan struct{}, _ *websocket.Conn)) Option {
	return func(c *Client) { c.pingLoop = f }
}

func WithSendOnConnect(f func() (int, []byte)) Option {
	return func(c *Client) { c.onConnect = append(c.onConnect, f) }
}

func WithContext(ctx context.Context) Option {
	return func(c *Client) {
		c.ctx, c.cancel = context.WithCancel(ctx)
	}
}

func WithCompressor(compressor compressor.Compressor) Option {
	return func(c *Client) { c.c = compressor }
}

func WithDefaultCompress(ok bool) Option {
	return func(c *Client) { c.enableCompress = ok }
}
