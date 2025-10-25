package gemini

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/CrazyThursdayV50/pkgo/file"
	defaultlogger "github.com/CrazyThursdayV50/pkgo/log/default"
	"google.golang.org/genai"
)

func TestChat(t *testing.T) {
	var config Config
	config.Model = "gemini-2.5-flash"
	config.SystemFile = "../.system"
	config.SetThinkingConfig(&genai.ThinkingConfig{
		IncludeThoughts: false,
	})
	ctx := context.TODO()
	logger := defaultlogger.New(defaultlogger.DefaultConfig())
	logger.Init()

	q := "who are you?"

	client, err := New(ctx, logger, &config)
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

func TestChatStream(t *testing.T) {
	var config Config
	config.Model = "gemini-2.5-pro"
	config.SystemFile = "../.system"
	config.SetThinkingConfig(&genai.ThinkingConfig{
		IncludeThoughts: true,
	})
	ctx := context.TODO()
	logger := defaultlogger.New(defaultlogger.DefaultConfig())
	logger.Init()

	q, _ := file.ReadFileToString("../question.json")
	// q := "hello, who are you?"

	client, err := New(ctx, logger, &config)
	if err != nil {
		t.Fatalf("new client failed: %v", err)
		return
	}

	start := time.Now()
	textChan, errChan := client.ChatStream(ctx, q)
	cost := time.Since(start)
	fmt.Printf("Cost: %s\n", cost)

	go func(t *testing.T) {
		for err := range errChan {
			fmt.Printf("chat stream failed: %v\n", err)
		}
	}(t)

	for text := range textChan {
		fmt.Print(text)
	}
	println()
}
