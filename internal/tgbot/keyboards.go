package tgbot

import (
	"fmt"
	"time"

	"github.com/go-telegram/bot/models"
	"github.com/nerthisdev/lab-q/internal/repository"
)

type Keyboards struct{}

func buildMainMenu(isAdmin bool) *models.ReplyKeyboardMarkup {
	rows := [][]models.KeyboardButton{{
		{Text: "Join Queue"}, {Text: "Check Queue"}, {Text: "List Subjects"}, {Text: "Change Name"},
	}}
	if isAdmin {
		rows = append(rows, []models.KeyboardButton{{Text: "Add Class"}, {Text: "Add Date"}})
	}
	return &models.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}
}

func buildSubjectsInline(subjects []repository.Subject) *models.InlineKeyboardMarkup {
	var rows [][]models.InlineKeyboardButton
	for _, s := range subjects {
		btn := models.InlineKeyboardButton{Text: s.Name, CallbackData: fmt.Sprintf("show_subj_%d", s.ID)}
		rows = append(rows, []models.InlineKeyboardButton{btn})
	}
	return &models.InlineKeyboardMarkup{InlineKeyboard: rows}
}

func buildDatesInline(schedules []repository.SubjectSchedule) *models.InlineKeyboardMarkup {
	var rows [][]models.InlineKeyboardButton
	now := time.Now()
	for _, sch := range schedules {
		t := nextOccurrence(sch, now)
		btnText := t.Format("2006-01-02 15:04")
		cbData := fmt.Sprintf("show_date_%d_%d", sch.ID, t.Unix())
		rows = append(rows, []models.InlineKeyboardButton{{Text: btnText, CallbackData: cbData}})
	}
	if len(rows) == 0 {
		rows = append(rows, []models.InlineKeyboardButton{{Text: "no dates", CallbackData: "noop"}})
	}
	return &models.InlineKeyboardMarkup{InlineKeyboard: rows}
}

func nextOccurrence(sch repository.SubjectSchedule, after time.Time) time.Time {
	t := time.Date(sch.StartDate.Year(), sch.StartDate.Month(), sch.StartDate.Day(),
		sch.TimeOfDay.Hour(), sch.TimeOfDay.Minute(), 0, 0, after.Location())
	for t.Before(after) {
		t = t.AddDate(0, 0, 7*sch.IntervalWeeks)
	}
	return t
}
