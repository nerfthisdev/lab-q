package tgbot

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

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

// StartHandler registers a user and sends welcome message
func (tgb *Tgbot) StartHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	u := update.Message.From
	user := repository.User{
		TelegramUserID: u.ID,
		TelegramChatID: update.Message.Chat.ID,
		Username:       strings.TrimSpace(u.Username),
		IsAdmin:        false,
	}
	if err := tgb.Repository.CreateOrUpdateUser(user); err != nil {
		tgb.Logger.Error("failed to register user", zap.Error(err))
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Welcome to LabQ bot!",
	})
}

// AddClassHandler adds a new subject (admin only)
func (tgb *Tgbot) AddClassHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	from := update.Message.From
	u, err := tgb.Repository.GetUserByID(from.ID)
	if err != nil || !u.IsAdmin {
		return
	}

	args := strings.Fields(update.Message.Text)
	if len(args) < 2 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "usage: /add_class <name>",
		})
		return
	}

	name := strings.Join(args[1:], " ")
	if err := tgb.Repository.CreateSubject(name, ""); err != nil {
		tgb.Logger.Error("failed to create subject", zap.Error(err))
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "failed to add class"})
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "class added"})
}

// AddDateHandler adds schedule entry (admin only)
func (tgb *Tgbot) AddDateHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	from := update.Message.From
	u, err := tgb.Repository.GetUserByID(from.ID)
	if err != nil || !u.IsAdmin {
		return
	}

	args := strings.Fields(update.Message.Text)
	if len(args) < 6 {
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "usage: /add_date <subject_id> <day_of_week> <time(HH:MM)> <interval_weeks> <start_date(YYYY-MM-DD)>"})
		return
	}

	subjectID := args[1]
	day := args[2]
	timeStr := args[3]
	interval := args[4]
	startDateStr := args[5]

	sid, err := strconv.ParseInt(subjectID, 10, 64)
	if err != nil {
		return
	}
	dow, err := strconv.Atoi(day)
	if err != nil {
		return
	}
	tod, err := time.Parse("15:04", timeStr)
	if err != nil {
		return
	}
	iv, err := strconv.Atoi(interval)
	if err != nil {
		return
	}
	start, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return
	}

	if err := tgb.Repository.CreateSubjectSchedule(sid, dow, tod, start, iv); err != nil {
		tgb.Logger.Error("failed to create schedule", zap.Error(err))
		return
	}
	b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "date added"})
}

// JoinHandler allows a user to join the queue
func (tgb *Tgbot) JoinHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	args := strings.Fields(update.Message.Text)
	if len(args) < 2 {
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "usage: /join <subject_id>"})
		return
	}

	sid, err := strconv.Atoi(args[1])
	if err != nil {
		return
	}

	if err := tgb.Repository.AddUserToQueue(sid, update.Message.From.ID); err != nil {
		tgb.Logger.Error("failed to add to queue", zap.Error(err))
		return
	}
	b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "joined queue"})
}

// QueueHandler prints the queue for a subject
func (tgb *Tgbot) QueueHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	args := strings.Fields(update.Message.Text)
	if len(args) < 2 {
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "usage: /queue <subject_id>"})
		return
	}

	sid, err := strconv.Atoi(args[1])
	if err != nil {
		return
	}

	users, err := tgb.Repository.GetQueueForSubject(sid)
	if err != nil {
		tgb.Logger.Error("failed to get queue", zap.Error(err))
		return
	}

	var bld strings.Builder
	for i, u := range users {
		bld.WriteString(fmt.Sprintf("%d. %s\n", i+1, u.Username))
	}
	if bld.Len() == 0 {
		bld.WriteString("queue is empty")
	}
	b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: bld.String()})
}
