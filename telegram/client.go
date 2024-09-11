package telegram

import (
	"context"

	"github.com/CrazyThursdayV50/gotils/pkg/async/goo"
	"github.com/CrazyThursdayV50/pkgo/trace"
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	KYCGroupID int64
	*tgbot.BotAPI
	tracer trace.Tracer
}

type UpdateHandler func(context.Context, trace.Tracer, tgbotapi.Update)

func New(cfg *Config, tracer trace.Tracer) (*Bot, error) {
	bot, err := tgbot.NewBotAPI(cfg.APIKEY)
	if err != nil {
		return nil, err
	}
	bot.Debug = cfg.Debug
	return &Bot{KYCGroupID: cfg.KycGroupID, BotAPI: bot, tracer: tracer}, nil
}

func (b *Bot) Run(ctx context.Context, handler UpdateHandler) error {
	span, ctx := b.tracer.NewSpan(ctx)
	defer span.Finish()
	ch := b.BotAPI.GetUpdatesChan(tgbot.UpdateConfig{})
	goo.Go(func() {
		for update := range ch {
			select {
			case <-ctx.Done():
			default:
				handler(ctx, b.tracer, update)
			}
		}
	})
	return nil
}
