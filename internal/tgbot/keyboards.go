package tgbot

import "github.com/go-telegram/bot/models"

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
