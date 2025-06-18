package main

import (
	"context"
	"flag"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/nerthisdev/lab-q/internal/config"
	"github.com/nerthisdev/lab-q/internal/database"
	"github.com/nerthisdev/lab-q/internal/logger"
	"github.com/nerthisdev/lab-q/internal/repository"
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

	db, err := database.NewPostgresDB(cfg.Database)
	if err != nil {
		logger.Fatal("failed to connect to DB", zap.String("reason", err.Error()))
	}
	repo := repository.Init(db)

	defer cancel()

	opts := []bot.Option{}

	tgb := tgbot.Init(cfg, opts, ctx, &logger, repo)
	tgb.Run(ctx)
}
