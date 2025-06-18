package tgbot

import (
	"context"
	"log"
	"net/http"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/nerthisdev/lab-q/internal/config"
	"go.uber.org/zap"
)

type Tgbot struct {
	Bot    *bot.Bot
	Logger *zap.Logger
	Config *config.Config
}

func Init(config *config.Config, opts []bot.Option, ctx context.Context, logger *zap.Logger) *Tgbot {
	tgb := Tgbot{
		Logger: logger,
		Config: config,
	}

	b, err := bot.New(tgb.Config.Bot.APIToken, opts...)
	if err != nil {
		tgb.Logger.Fatal("failed to create new bot instance", zap.String("error", err.Error()))
	}
	tgb.Logger.Info("created new bot instance")

	tgb.Bot = b

	button := bot.SetChatMenuButtonParams{
		MenuButton: models.MenuButtonCommands{},
	}

	tgb.Bot.SetChatMenuButton(ctx, &button)

	set, err := tgb.Bot.SetWebhook(ctx, &bot.SetWebhookParams{
		URL: tgb.Config.Bot.WebHookURL,
	})
	if err != nil {
		tgb.Logger.Fatal("failed to set webhook", zap.String("error", err.Error()))
		log.Fatal(err.Error())
	}

	if !set {
		return nil
	}

	return &tgb
}

func (tgb *Tgbot) Run(ctx context.Context) {
	go func() {
		http.ListenAndServe(tgb.Config.Server.Port, tgb.Bot.WebhookHandler())
	}()

	tgb.Bot.StartWebhook(ctx)
}

func (tgb *Tgbot) registerCommands(ctx context.Context, commands []models.BotCommand) {
	tgb.Bot.SetMyCommands(ctx, &bot.SetMyCommandsParams{
		Commands: commands,
	})
}
