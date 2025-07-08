package tgbot

import (
    "context"

    "github.com/go-telegram/bot"
    "github.com/go-telegram/bot/models"
    "go.uber.org/zap"
)

// logUpdateMiddleware logs information about every incoming update.
func (tgb *Tgbot) logUpdateMiddleware(next bot.HandlerFunc) bot.HandlerFunc {
    return func(ctx context.Context, b *bot.Bot, update *models.Update) {
        if update.Message != nil {
            from := update.Message.From
            tgb.Logger.Info("incoming message",
                zap.Int64("telegram_id", from.ID),
                zap.String("username", from.Username),
                zap.String("text", update.Message.Text),
            )
        }
        next(ctx, b, update)
    }
}
