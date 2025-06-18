package tgbot

import (
	"context"
	"fmt"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	m, errSend := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   update.Message.Text,
	})
	if errSend != nil {
		fmt.Printf("error sending message: %v\n", errSend)
		return
	}

	time.Sleep(time.Second * 2)

	_, errEdit := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:    m.Chat.ID,
		MessageID: m.ID,
		Text:      "New Message!",
	})
	if errEdit != nil {
		fmt.Printf("error edit message: %v\n", errEdit)
		return
	}
}
