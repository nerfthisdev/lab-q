package bot

import (
	"context"
	"net/http"

	"github.com/go-telegram/bot"
	"go.uber.org/zap"
)

type Tgbot struct {
	Bot    *bot.Bot
	Logger *zap.Logger
	port   string
	ctx    context.Context
}

func Init(apiKey string, opts []bot.Option, webHookURL string, ctx context.Context, logger zap.Logger, port string) *Tgbot {
	tgb := Tgbot{
		Logger: &logger,
	}

	b, err := bot.New(apiKey, opts...)
	if err != nil {
		tgb.Logger.Fatal("failed to create new bot instance", zap.String("error", err.Error()))
	}

	tgb.Bot = b

	tgb.Bot.SetWebhook(ctx, &bot.SetWebhookParams{
		URL: webHookURL,
	})

	return &tgb
}

func (tgb *Tgbot) Run() {
	go func() {
		http.ListenAndServe(tgb.port, tgb.Bot.WebhookHandler())
	}()

	tgb.Bot.StartWebhook(tgb.ctx)
}
