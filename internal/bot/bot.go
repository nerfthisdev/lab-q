package bot

import (
	"github.com/go-telegram/bot"
)

type tgbot struct{}

func init(apiKey string, opts []bot.Option) {
}

func createBotInstance(apiKey string, opts []bot.Option) (bot.Bot, error) {
	b, err := bot.New(apiKey, opts...)
	if err != nil {
		return nil, err
	}

	return b, nil
}
