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

	// If we previously asked for a name, treat the incoming text as the user's name
	if tgb.awaitingName[userID] {
		if text == "" {
			b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: "name cannot be empty"})
			return
		}
		err := tgb.Repository.CreateOrUpdateUser(repository.User{
			TelegramUserID: userID,
			TelegramChatID: chatID,
			Username:       text,
			IsAdmin:        tgb.isAdmin(userID),
		})
		if err != nil {
			tgb.Logger.Error("could not save user name", zap.Error(err))
			b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: "failed to save name"})
			return
		}
		delete(tgb.awaitingName, userID)
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: "Name saved"})
		return
	}

	// Unknown text message
	switch strings.ToLower(text) {
	case "join queue":
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: "use /join <subject_id>"})
	case "check queue":
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: "use /queue <subject_id>"})
	case "list subjects":
		tgb.SubjectsHandler(ctx, b, update)
	case "change name":
		tgb.SetNameHandler(ctx, b, update)
	case "add class":
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: "use /add_class <name>"})
	case "add date":
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: "use /add_date <subject_id> <day_of_week> <time(HH:MM)> <interval_weeks> <start_date(YYYY-MM-DD)>"})
	default:
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: "unknown command"})
	}
}

// StartHandler registers a user and sends welcome message
func (tgb *Tgbot) StartHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	u := update.Message.From
	chatID := update.Message.Chat.ID

	dbUser, err := tgb.Repository.GetUserByID(u.ID)
	if err != nil {
		// new user
		user := repository.User{
			TelegramUserID: u.ID,
			TelegramChatID: chatID,
			Username:       "",
			IsAdmin:        tgb.isAdmin(u.ID),
		}
		if err := tgb.Repository.CreateOrUpdateUser(user); err != nil {
			tgb.Logger.Error("failed to register user", zap.Error(err))
			return
		}
		tgb.awaitingName[u.ID] = true
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: "Welcome! Please send your name."})
		return
	}

	// existing user - update chat id
	dbUser.TelegramChatID = chatID
	if err := tgb.Repository.CreateOrUpdateUser(*dbUser); err != nil {
		tgb.Logger.Error("failed to update user", zap.Error(err))
	}

	if dbUser.Username == "" {
		tgb.awaitingName[u.ID] = true
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: "Please send your name."})
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        fmt.Sprintf("Welcome back, %s!", dbUser.Username),
		ReplyMarkup: buildMainMenu(dbUser.IsAdmin),
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

// SubjectsHandler lists all available subjects
func (tgb *Tgbot) SubjectsHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	subjects, err := tgb.Repository.GetAllSubjects()
	if err != nil {
		tgb.Logger.Error("failed to list subjects", zap.Error(err))
		return
	}

	if len(subjects) == 0 {
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "no subjects"})
		return
	}

	var bld strings.Builder
	for _, s := range subjects {
		bld.WriteString(fmt.Sprintf("%d. %s\n", s.ID, s.Name))
	}
	b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: bld.String()})
}

// SetNameHandler prompts user to enter a new name
func (tgb *Tgbot) SetNameHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	tgb.awaitingName[update.Message.From.ID] = true
	b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "Please send your new name."})
}

// helper to check admin status
func (tgb *Tgbot) isAdmin(id int64) bool {
	for _, a := range tgb.Config.AdminIDs {
		if a == id {
			return true
		}
	}
	return false
}

// ShowHandler presents subject list with inline buttons
func (tgb *Tgbot) ShowHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	subjects, err := tgb.Repository.GetAllSubjects()
	if err != nil {
		tgb.Logger.Error("failed to list subjects", zap.Error(err))
		return
	}

	kb := buildSubjectsInline(subjects)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        "Choose subject:",
		ReplyMarkup: kb,
	})
}

// ShowSubjectCallback displays schedule dates for the chosen subject
func (tgb *Tgbot) ShowSubjectCallback(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.CallbackQuery == nil {
		return
	}

	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID})

	data := strings.TrimPrefix(update.CallbackQuery.Data, "show_subj_")
	sid, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		return
	}

	schedules, err := tgb.Repository.GetSchedulesForSubject(sid)
	if err != nil {
		tgb.Logger.Error("failed to get schedules", zap.Error(err))
		return
	}

	kb := buildDatesInline(schedules)
	b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		MessageID:   update.CallbackQuery.Message.Message.ID,
		Text:        "Choose date:",
		ReplyMarkup: kb,
	})
}

// ShowDateCallback finalizes date selection
func (tgb *Tgbot) ShowDateCallback(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.CallbackQuery == nil {
		return
	}

	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID})

	parts := strings.Split(strings.TrimPrefix(update.CallbackQuery.Data, "show_date_"), "_")
	if len(parts) != 2 {
		return
	}

	ts, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return
	}

	dt := time.Unix(ts, 0)
	b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
		MessageID: update.CallbackQuery.Message.Message.ID,
		Text:      fmt.Sprintf("Selected date: %s", dt.Format("2006-01-02 15:04")),
	})
}
