package server

import (
	"context"

	"github.com/CrazyThursdayV50/pkgo/trace"
	"github.com/opentracing/opentracing-go/log"
)

func (c *conn) send(ctx context.Context, tracer trace.Tracer, messageType int, data []byte) error {
	span, ctx := tracer.NewSpan(ctx)
	defer span.Finish()

	err := c.conn.WriteMessage(messageType, data)
	if err != nil {
		span.LogFields(log.Event("WriteMessage"), log.Error(err))
		c.cancel()
		return err
	}

	return nil
}

func (c *conn) recv(ctx context.Context, tracer trace.Tracer) (int, []byte, error) {
	span, ctx := tracer.NewSpan(ctx)
	defer span.Finish()
	messageType, data, err := c.conn.ReadMessage()
	if err != nil {
		span.LogFields(log.Event("ReadMessage"), log.Error(err))
		c.cancel()
		return messageType, nil, err
	}

	return messageType, data, nil
}
