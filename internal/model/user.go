package model

import "time"

type UserProfile struct {
	TelegramUserID int64
	TelegramChatID int64
	CreatedAt      time.Time
}
