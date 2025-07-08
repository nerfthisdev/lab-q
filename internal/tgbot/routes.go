package tgbot

import "github.com/go-telegram/bot"

func (tgb *Tgbot) RegisterRoutes() {
	tgb.Bot.RegisterHandler(bot.HandlerTypeMessageText, "start", bot.MatchTypeExact, tgb.StartHandler)
	tgb.Bot.RegisterHandler(bot.HandlerTypeMessageText, "add_class", bot.MatchTypeCommand, tgb.AddClassHandler)
	tgb.Bot.RegisterHandler(bot.HandlerTypeMessageText, "add_date", bot.MatchTypeCommand, tgb.AddDateHandler)
	tgb.Bot.RegisterHandler(bot.HandlerTypeMessageText, "join", bot.MatchTypeCommand, tgb.JoinHandler)
	tgb.Bot.RegisterHandler(bot.HandlerTypeMessageText, "queue", bot.MatchTypeCommand, tgb.QueueHandler)
	tgb.Bot.RegisterHandler(bot.HandlerTypeMessageText, "subjects", bot.MatchTypeCommand, tgb.SubjectsHandler)
	tgb.Bot.RegisterHandler(bot.HandlerTypeMessageText, "setname", bot.MatchTypeCommand, tgb.SetNameHandler)
	tgb.Bot.RegisterHandler(bot.HandlerTypeMessageText, "show", bot.MatchTypeCommand, tgb.ShowHandler)
	tgb.Bot.RegisterHandler(bot.HandlerTypeCallbackQueryData, "show_subj_", bot.MatchTypePrefix, tgb.ShowSubjectCallback)
	tgb.Bot.RegisterHandler(bot.HandlerTypeCallbackQueryData, "show_date_", bot.MatchTypePrefix, tgb.ShowDateCallback)
}
