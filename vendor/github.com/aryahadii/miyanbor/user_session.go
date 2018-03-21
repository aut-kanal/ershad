package miyanbor

import telegramAPI "gopkg.in/telegram-bot-api.v4"

type UserSession struct {
	User            *telegramAPI.User
	UserID          int
	ChatID          int64
	Payload         map[string]interface{}
	messageCallback CallbackFunction
}
