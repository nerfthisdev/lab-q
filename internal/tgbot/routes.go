package tgbot

import "github.com/go-telegram/bot"

func (tgb *Tgbot) registerRoutes() {
	tgb.Bot.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, tgb.StartHandler)
	tgb.Bot.RegisterHandler(bot.HandlerTypeMessageText, "/add_class", bot.MatchTypeCommand, tgb.AddClassHandler)
	tgb.Bot.RegisterHandler(bot.HandlerTypeMessageText, "/add_date", bot.MatchTypeCommand, tgb.AddDateHandler)
	tgb.Bot.RegisterHandler(bot.HandlerTypeMessageText, "/join", bot.MatchTypeCommand, tgb.JoinHandler)
	tgb.Bot.RegisterHandler(bot.HandlerTypeMessageText, "/queue", bot.MatchTypeCommand, tgb.QueueHandler)
}
