package tgbot

import (
	"context"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/nerthisdev/lab-q/internal/repository"
	"go.uber.org/zap"
)

func (tgb *Tgbot) DefaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	from := update.Message.From
	userID := from.ID
	chatID := update.Message.Chat.ID
	text := strings.TrimSpace(update.Message.Text)

	name := text
	err := tgb.Repository.CreateOrUpdateUser(repository.User{
		TelegramUserID: userID,
		TelegramChatID: chatID,
		Username:       name,
		IsAdmin:        false,
	})
	if err != nil {
		tgb.Logger.Fatal("couldnt create or update user", zap.String("reason", err.Error()))
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Error occured while regestering",
		})
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "Welcome",
	})

}
