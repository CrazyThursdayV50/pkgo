package resty

import (
	"context"
	"time"

	"github.com/CrazyThursdayV50/pkgo/json"
	"github.com/CrazyThursdayV50/pkgo/log"
	"github.com/CrazyThursdayV50/pkgo/request/resty/callbacks"
	"github.com/CrazyThursdayV50/pkgo/trace"
	"github.com/go-resty/resty/v2"
)

type Client struct {
	tracer trace.Tracer
	cfg    *Config
	client *resty.Client
	logger log.Logger
}

type RestyClient = resty.Client
type RestyRequest = resty.Request

func (c *Client) Request(ctx context.Context) *resty.Request {
	return c.client.Clone().R().SetContext(ctx)
}

type Option func(*Client)

func WithLogger(l log.Logger) Option {
	return func(c *Client) {
		c.logger = l
		if c.client != nil {
			c.client.SetLogger(l)
		}
	}
}

func WithOnBeforeRequest(f func(*RestyClient, *RestyRequest) error) Option {
	return func(c *Client) {
		c.client.OnBeforeRequest(f)
	}
}

func WithConfig(conf *Config) Option {
	return func(cl *Client) {
		cl.cfg = conf
	}
}

func WithTracer(tracer trace.Tracer) Option {
	return func(cl *Client) {
		cl.tracer = tracer
	}
}

func New(opts ...Option) *Client {
	var c Client

	for _, opt := range opts {
		opt(&c)
	}

	client := resty.New()
	client.SetDebug(c.cfg.Debug)
	client.SetTimeout(c.cfg.Timeout * time.Second)
	client.SetRetryCount(c.cfg.RetryCount)
	client.SetCloseConnection(c.cfg.CloseConnection)
	client.SetRedirectPolicy(resty.NoRedirectPolicy())
	if c.cfg.RetryWaitTime != 0 {
		client.SetRetryWaitTime(c.cfg.RetryWaitTime * time.Second)
	}
	if c.cfg.RetryMaxWaitTime != 0 {
		client.SetRetryMaxWaitTime(c.cfg.RetryMaxWaitTime * time.Second)
	}

	client.JSONMarshal = json.Marshal
	client.JSONUnmarshal = json.Unmarshal
	client.SetLogger(c.logger)

	c.client = client

	if c.cfg.EnableTrace {
		c.client.OnBeforeRequest(callbacks.TraceRequest(c.tracer))
		c.client.OnAfterResponse(callbacks.TraceResponse(c.tracer))
	}

	if c.cfg.EnableLog {
		c.client.OnBeforeRequest(callbacks.LogBeforeRequest(c.logger))
		c.client.OnAfterResponse(callbacks.LogAfterResponse(c.logger))
		c.client.OnError(callbacks.LogOnError(c.logger))
	}

	return &c
}
