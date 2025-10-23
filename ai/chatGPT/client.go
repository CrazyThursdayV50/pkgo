package chatgpt

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/CrazyThursdayV50/pkgo/log"
	gpt3 "github.com/sashabaranov/go-openai"
)

type Client struct {
	cfg    *Config
	client *gpt3.Client
	system string
}

var (
	ChatGPT4     = gpt3.GPT4Turbo
	ChatGPT4o    = gpt3.GPT4o
	ChatGPT3Dot5 = gpt3.GPT3Dot5Turbo
)

func New(cfg *Config, logger log.Logger) *Client {
	if cfg.Token == "" {
		cfg.Token = os.Getenv("OPENAI_API_KEY")
	}

	client := gpt3.NewClient(cfg.Token)

	var systemData []byte
	if cfg.SystemFile != "" {
		file, err := os.Open(cfg.SystemFile)
		if err != nil {
			logger.DPanicf("open system file failed: %v", err)
		}
		defer file.Close()

		systemData, err = io.ReadAll(file)
		if err != nil {
			logger.DPanicf("read system file failed: %v", err)
		}
	}

	return &Client{cfg: cfg, client: client, system: string(systemData)}
}

func (c *Client) Chat(ctx context.Context, q string, model string) (string, error) {
	resp, err := c.client.CreateChatCompletion(ctx, gpt3.ChatCompletionRequest{
		Model:               model,
		MaxCompletionTokens: c.cfg.MaxTokens,
		Messages: []gpt3.ChatCompletionMessage{
			{
				Role:    gpt3.ChatMessageRoleSystem,
				Content: c.system,
			},
			{
				Role:    gpt3.ChatMessageRoleUser,
				Content: q,
			},
		},
		Temperature: c.cfg.Temperature,
	})
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}

func (c *Client) ChatStream(ctx context.Context, q string, model string) (string, error) {
	stream, err := c.client.CreateChatCompletionStream(ctx, gpt3.ChatCompletionRequest{
		Model:               model,
		MaxCompletionTokens: c.cfg.MaxTokens,
		Messages: []gpt3.ChatCompletionMessage{
			{
				Role:    gpt3.ChatMessageRoleSystem,
				Content: c.system,
			},
			{
				Role:    gpt3.ChatMessageRoleUser,
				Content: q,
			},
		},
		Temperature: c.cfg.Temperature,
	})
	if err != nil {
		return "", err
	}
	defer stream.Close()

	var answer string
	for {
		resp, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return answer, nil
		}
		if err != nil {
			return "", err
		}
		answer += resp.Choices[0].Delta.Content
	}
}

func (c *Client) System() string { return c.system }
