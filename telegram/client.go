package telegram

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/url"

	"github.com/CrazyThursdayV50/pkgo/goo"
	"github.com/CrazyThursdayV50/pkgo/log"
	"github.com/CrazyThursdayV50/pkgo/trace"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	*tgbotapi.BotAPI
	tracer       trace.Tracer
	logger       log.Logger
	updateFilter []string
}

type Update = tgbotapi.Update

type UpdateHandler func(context.Context, trace.Tracer, Update, *Bot)

func New(cfg *Config, logger log.Logger, tracer trace.Tracer) (*Bot, error) {
	var client http.Client

	var transport = http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: cfg.SkipTlsVerify},
	}

	if cfg.Proxy != "" {
		url, err := url.Parse(cfg.Proxy)
		if err != nil {
			return nil, err
		}

		transport.Proxy = http.ProxyURL(url)
	}

	client.Transport = &transport

	bot, err := tgbotapi.NewBotAPIWithClient(cfg.APIKEY, tgbotapi.APIEndpoint, &client)
	if err != nil {
		return nil, err
	}
	bot.Debug = cfg.Debug
	return &Bot{BotAPI: bot, tracer: tracer, logger: logger, updateFilter: cfg.UpdateFilter}, nil
}

func (b *Bot) Run(ctx context.Context, handler UpdateHandler) error {
	span, ctx := b.tracer.NewSpan(ctx)
	defer span.Finish()
	cfg := tgbotapi.NewUpdate(0)
	cfg.Timeout = 60
	cfg.AllowedUpdates = append(make([]string, 0), b.updateFilter...)
	ch := b.GetUpdatesChan(cfg)
	goo.Goo(func() {
		for update := range ch {
			select {
			case <-ctx.Done():
			default:
				handler(ctx, b.tracer, update, b)
			}
		}
	}, func(err error) { b.logger.Errorf("handler panic: %v", err) })
	return nil
}
