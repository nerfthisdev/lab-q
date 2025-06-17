package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/nerthisdev/lab-q/internal/logger"
	"github.com/nerthisdev/lab-q/internal/tgbot"
)

func main() {
	logger := logger.GetLogger()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	token := os.Getenv("tg_api")

	tgb := tgbot.Init(token, opts, "https://a3da-5-180-114-202.ngrok-free.app", ctx, &logger, ":2000")
	tgb.Run(ctx)
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   update.Message.Text,
	})
}
