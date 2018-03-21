package keyboard

import (
	"fmt"

	"gitlab.com/kanalbot/ershad/ui/text"
	telegramAPI "gopkg.in/telegram-bot-api.v4"
)

const (
	KeyboardAcceptButtonData = "a%s"
	KeyboardRejectButtonData = "r%s"
)

func NewErshadInlineKeyboard(id string) telegramAPI.InlineKeyboardMarkup {
	var row []telegramAPI.InlineKeyboardButton

	acceptData := fmt.Sprintf(KeyboardAcceptButtonData, id)
	acceptButton := telegramAPI.NewInlineKeyboardButtonData(text.KeyboardAcceptButton, acceptData)

	rejectData := fmt.Sprintf(KeyboardRejectButtonData, id)
	rejectButton := telegramAPI.NewInlineKeyboardButtonData(text.KeyboardRejectButton, rejectData)

	row = append(row, acceptButton, rejectButton)

	return telegramAPI.NewInlineKeyboardMarkup(row)
}
