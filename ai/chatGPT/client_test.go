package chatgpt

import (
	"context"
	"fmt"
	"testing"
	"time"

	defaultlogger "github.com/CrazyThursdayV50/pkgo/log/default"
)

func TestChat(t *testing.T) {
	var config Config
	config.Model = "gemini-2.5-flash"
	config.SystemFile = "../.system"
	ctx := context.TODO()
	logger := defaultlogger.New(defaultlogger.DefaultConfig())
	logger.Init()

	q := "who are you?"

	client, err := New(&config, logger)
	if err != nil {
		t.Fatalf("new client failed: %v", err)
		return
	}

	start := time.Now()
	text, err := client.Chat(ctx, q)
	cost := time.Since(start)
	if err != nil {
		t.Fatalf("chat failed: %v", err)
		return
	}
	fmt.Printf("[%s]receive: %s\n", cost.String(), text)

}
