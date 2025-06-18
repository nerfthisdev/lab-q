package main

import (
	"context"
	"flag"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/nerthisdev/lab-q/internal/config"
	"github.com/nerthisdev/lab-q/internal/logger"
	"github.com/nerthisdev/lab-q/internal/tgbot"
	"go.uber.org/zap"
)

func main() {
	configPath := flag.String("c", "./cmd/config.yaml", "path to go-telegram-bot-example config")

	flag.Parse()

	logger := logger.GetLogger()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	cfg := &config.Config{}

	err := config.GetConfiguration(*configPath, cfg)
	if err != nil {
		logger.Fatal("failed to get configuration", zap.String("reason", err.Error()))
	}

	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(tgbot.Handler),
	}

	tgb := tgbot.Init(cfg, opts, ctx, &logger)
	tgb.Run(ctx)
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   update.Message.Text,
	})
}
