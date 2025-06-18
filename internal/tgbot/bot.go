package tgbot

import (
	"context"
	"log"
	"net/http"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"
)

type Tgbot struct {
	Bot        *bot.Bot
	Logger     *zap.Logger
	port       string
	webHookURL string
}

func Init(apiKey string, opts []bot.Option, webHookURL string, ctx context.Context, logger *zap.Logger, port string) *Tgbot {
	commands := []models.BotCommand{
		{
			Command:     "start",
			Description: "launch the bot",
		},
		{
			Command:     "register",
			Description: "register in the bot",
		},
	}

	tgb := Tgbot{
		Logger:     logger,
		webHookURL: webHookURL,
	}

	b, err := bot.New(apiKey, opts...)
	if err != nil {
		tgb.Logger.Fatal("failed to create new bot instance", zap.String("error", err.Error()))
	}
	tgb.Logger.Info("created new bot instance")

	tgb.Bot = b
	tgb.registerCommands(ctx, commands)

	button := bot.SetChatMenuButtonParams{
		MenuButton: models.MenuButtonCommands{},
	}

	tgb.Bot.SetChatMenuButton(ctx, &button)

	set, err := tgb.Bot.SetWebhook(ctx, &bot.SetWebhookParams{
		URL: tgb.webHookURL,
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
		http.ListenAndServe(":2000", tgb.Bot.WebhookHandler())
	}()

	tgb.Bot.StartWebhook(ctx)
}

func (tgb *Tgbot) registerCommands(ctx context.Context, commands []models.BotCommand) {
	tgb.Bot.SetMyCommands(ctx, &bot.SetMyCommandsParams{
		Commands: commands,
	})

}
