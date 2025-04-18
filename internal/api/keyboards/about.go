package keyboards

import "github.com/go-telegram/bot/models"

var (
	About = models.InlineKeyboardMarkup{InlineKeyboard: [][]models.InlineKeyboardButton{
		{
			models.InlineKeyboardButton{
				Text:         "Назад",
				CallbackData: "start",
			},
		},
	}}
)
