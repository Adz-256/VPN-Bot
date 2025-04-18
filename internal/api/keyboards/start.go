package keyboards

import "github.com/go-telegram/bot/models"

var (
	Start = models.InlineKeyboardMarkup{InlineKeyboard: [][]models.InlineKeyboardButton{
		{
			models.InlineKeyboardButton{
				Text:         "Купить Подписку",
				CallbackData: "buy",
			},
		},
		{
			models.InlineKeyboardButton{
				Text:         "Мои Подписки",
				CallbackData: "subscriptions",
			},
		},
		{
			models.InlineKeyboardButton{
				Text:         "Поддержка",
				CallbackData: "support",
			},
			models.InlineKeyboardButton{
				Text:         "Инструкция",
				CallbackData: "instructions",
			},
		},
	},
	}
)
