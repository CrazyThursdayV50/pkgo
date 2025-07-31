package websocket

import (
	"context"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/CrazyThursdayV50/pkgo/goo"
	"github.com/CrazyThursdayV50/pkgo/log"
	defaultlogger "github.com/CrazyThursdayV50/pkgo/log/default"
	"github.com/CrazyThursdayV50/pkgo/trace/jaeger"
	client "github.com/CrazyThursdayV50/pkgo/websocket/client"
	"github.com/CrazyThursdayV50/pkgo/websocket/server"
	"github.com/gorilla/websocket"
)

func TestWebsocket(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := defaultlogger.DefaultConfig()
	cfg.Level = "info"
	var logger = defaultlogger.New(cfg)
	logger.Init()

	cfg.Level = "fatal"
	var jaegerLogger = defaultlogger.New(cfg)
	jaegerLogger.Init()

	jaegerCfg := jaeger.DefaultConfig()
	jaegerCfg.LogSpans = false
	tracer, err := jaeger.New(ctx, jaegerCfg, jaegerLogger)
	if err != nil {
		t.Fatalf("logger failed: %v", err)
	}

	var start = time.Now()
	var count int64

	wsserver := server.New(
		server.WithLogger(logger),
		server.WithTracer(tracer.NewTracer("WsServer")),
		server.WithHandler(func(ctx context.Context, messageType int, data []byte, err error) (int, []byte, error) {
			if err != nil {
				logger.Errorf("receive error: %v", err)
				return 0, nil, err
			}

			// logger.Infof("server receive: %s\n", data)

			if count == 0 {
				start = time.Now()
			}

			atomic.AddInt64(&count, 1)

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

	// goo.Goo(func() {
	// 	for {
	// 		wsserver.Broadcast(ctx, websocket.TextMessage, []byte("broadcast"))
	// 		time.Sleep(time.Second)
	// 	}
	// }, func(err error) { logger.Error(err) })

	var clientCount int64 = 1000
	var messagePerClient int64 = 100
	var message [1024]byte

	// ---------- client
	// url := "ws://127.0.0.1:18080/ws"
	var newClient = func() *client.Client {
		var wg sync.WaitGroup
		wg.Add(1)
		wsclient := client.New(
			client.WithURL("ws://localhost:18080"),
			// client.WithURL(url),
			client.WithContext(ctx), client.WithLogger(logger),

			client.WithMessageHandler(func(ctx context.Context, l log.Logger, typ int, data []byte, f func(error)) (int, []byte) {
				l.Infof("client receive: %s", data)
				return websocket.BinaryMessage, nil
			}),

			client.WithSendOnConnect(func() (int, []byte) {
				wg.Done()
				return client.PingMessage, []byte{}
			}),

			// client.WithPingLoop(func(done <-chan struct{}, conn *websocket.Conn) {
			// 	for {
			// 		select {
			// 		case <-done:
			// 			return
			// 		default:
			// 			conn.WriteMessage(websocket.TextMessage, []byte("ping from go websocket client"))
			// 			time.Sleep(time.Second * 30)
			// 		}
			// 	}
			// }),
		)

		wsclient.Run()
		wg.Wait()
		return wsclient
	}

	// _ = newClient
	logger.Infof("start clients ...")
	var clients = []*client.Client{}
	for range clientCount {
		clients = append(clients, newClient())
	}
	logger.Infof("start clients finished")

	for _, c := range clients {
		go func(client *client.Client) {
			for range messagePerClient {
				client.Send(message[:])
			}
		}(c)
	}

	var totalMessages = messagePerClient * clientCount
	var ticker = time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			logger.Infof("%d/%d", count, totalMessages)

		default:
			if count == totalMessages {
				cost := time.Since(start)
				ability := time.Duration(totalMessages) * time.Second / cost
				logger.Infof("total: %d, cost: %s, ability: %d", totalMessages, cost, ability)
				time.Sleep(time.Second * 10)
				return
			}
		}
	}

}
