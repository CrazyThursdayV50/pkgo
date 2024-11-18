package websocket

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/CrazyThursdayV50/gotils/pkg/async/goo"
	"github.com/CrazyThursdayV50/pkgo/log"
	defaultlogger "github.com/CrazyThursdayV50/pkgo/log/default"
	"github.com/CrazyThursdayV50/pkgo/trace/jaeger"
	"github.com/CrazyThursdayV50/pkgo/websocket/client"
	"github.com/CrazyThursdayV50/pkgo/websocket/server"
	"github.com/gorilla/websocket"
)

func TestWebsocket(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var logger = defaultlogger.New(defaultlogger.DefaultConfig())
	logger.Init()

	tracer, err := jaeger.New(ctx, jaeger.DefaultConfig(), logger)
	if err != nil {
		t.Fatalf("logger failed: %v", err)
	}

	wsserver := server.New(
		server.WithLogger(logger),
		server.WithTracer(tracer.NewTracer("WsServer")),
		server.WithHandler(func(messageType int, data []byte, err error) (int, []byte, error) {
			if err != nil {
				logger.Errorf("receive error: %v", err)
				return 0, nil, err
			}

			switch messageType {
			case websocket.TextMessage:
				logger.Infof("server receive: %s\n", data)
				return messageType, data, nil

			default:
				return 0, nil, nil
			}
		}),
	)

	mux := http.NewServeMux()
	mux.Handle("/ws", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = wsserver.Run(ctx, w, r, nil)
	}))

	goo.Go(func() {
		http.ListenAndServe(":18080", mux)
	})

	goo.Goo(func() {
		for {
			wsserver.Broadcast(ctx, websocket.TextMessage, []byte("broadcast"))
			time.Sleep(time.Second)
		}
	}, func(err error) { logger.Error(err) })

	// ---------- client
	wsclient := client.New(
		client.WithURL("ws://localhost:18080/ws"),
		client.WithContext(ctx), client.WithLogger(logger),

		client.WithMessageHandler(func(ctx context.Context, l log.Logger, data []byte, f func(error)) []byte {
			l.Infof("client receive: %s", data)
			return nil
		}),

		client.WithPingLoop(func(done <-chan struct{}, conn *websocket.Conn) {
			for {
				select {
				case <-done:
					return
				default:
					conn.WriteMessage(websocket.TextMessage, []byte("ping"))
					time.Sleep(time.Second * 10)
				}
			}
		}),
	)

	wsclient.Run()

	<-make(chan struct{})
}
