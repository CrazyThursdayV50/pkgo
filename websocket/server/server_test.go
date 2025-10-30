package server

import (
	"context"
	"net/http"
	"testing"

	"github.com/CrazyThursdayV50/pkgo/goo"
	"github.com/CrazyThursdayV50/pkgo/log/sugar"
	"github.com/CrazyThursdayV50/pkgo/trace/jaeger"
	"github.com/gorilla/websocket"
)

func TestServer(t *testing.T) {
	ctx := context.TODO()

	cfg := sugar.DefaultConfig()
	cfg.Level = "debug"
	var logger = sugar.New(cfg)

	cfg.Level = "fatal"
	var jaegerLogger = sugar.New(cfg)

	jaegerCfg := jaeger.DefaultConfig()
	jaegerCfg.LogSpans = false
	tracer, err := jaeger.New(ctx, jaegerCfg, jaegerLogger)

	if err != nil {
		t.Fatalf("logger failed: %v", err)
	}
	wsserver := New(
		WithLogger(logger),
		WithTracer(tracer.NewTracer("WsServer")),
		WithHandler(func(ctx context.Context, messageType int, data []byte, err error) (int, []byte, error) {
			if err != nil {
				logger.Errorf("receive error: %v", err)
				return 0, nil, err
			}

			logger.Infof("server receive: %s\n", data)

			switch messageType {
			case websocket.TextMessage:
				// logger.Infof("server receive: %s\n", data)
				return messageType, nil, nil

			default:
				// logger.Infof("server receive[%d]: %v\n", messageType, data)
				return 0, nil, nil
			}
		}),
	)

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = wsserver.Run(ctx, w, r, nil)
	}))

	goo.Go(func() {
		http.ListenAndServe(":18080", mux)
	})

	<-make(chan struct{})
}
