package keyboards

import (
	"strings"

	"github.com/go-telegram/bot/models"
)

var (
	BuyChooseServer = models.InlineKeyboardMarkup{InlineKeyboard: [][]models.InlineKeyboardButton{
		{
			models.InlineKeyboardButton{
				Text:         "🇳🇱 Нидерланды",
				CallbackData: "buy_RU",
			},
		},
		{
			models.InlineKeyboardButton{
				Text:         "🇫🇮 Финляндия",
				CallbackData: "buy_FIN",
			},
		},
		{
			models.InlineKeyboardButton{
				Text:         "Назад",
				CallbackData: "start",
			},
		},
	}}
	Payment = models.InlineKeyboardMarkup{InlineKeyboard: [][]models.InlineKeyboardButton{}}
	// BuySuccess = models.InlineKeyboardMarkup{InlineKeyboard: [][]models.InlineKeyboardButton{
	// 	{
	// 		models.InlineKeyboardButton{
	// 			Text:         "QR-code",
	// 			CallbackData: "qr",
	// 		},
	// 		models.InlineKeyboardButton{
	// 			Text:         "Конфигурационный файл",
	// 			CallbackData: "config_",
	// 		},
	// 	},
	// 	{
	// 		models.InlineKeyboardButton{
	// 			Text:         "Инструкция",
	// 			CallbackData: "instructions",
	// 		},
	// 	},
	// 	{
	// 		models.InlineKeyboardButton{
	// 			Text:         "На главную",
	// 			CallbackData: "start",
	// 		},
	// 	},
	// }}
)

func PrePayment(data string, url string, transID string) models.InlineKeyboardMarkup {
	arrData := strings.Split(data, "_")
	return models.InlineKeyboardMarkup{InlineKeyboard: [][]models.InlineKeyboardButton{
		{
			models.InlineKeyboardButton{
				Text:         "🔗 Ссылка на оплату",
				URL:          url,
				CallbackData: "payment" + "_" + transID + "_" + strings.Join(arrData[1:len(arrData)-1], "_"),
			},
		},
		{
			models.InlineKeyboardButton{
				Text:         "✅ Подтверждаю",
				CallbackData: "confirm" + "_" + transID + "_" + strings.Join(arrData[1:len(arrData)-1], "_"),
			},
		},
		{
			models.InlineKeyboardButton{
				Text:         "Отмена",
				CallbackData: strings.Join(arrData[1:len(arrData)-1], "_"),
			},
		},
	}}
}

func BuySuccess(id string) models.InlineKeyboardMarkup {
	return models.InlineKeyboardMarkup{InlineKeyboard: [][]models.InlineKeyboardButton{
		{
			models.InlineKeyboardButton{
				Text:         "QR-code",
				CallbackData: "qr" + "_" + id,
			},
			models.InlineKeyboardButton{
				Text:         "Конфигурационный файл",
				CallbackData: "config_" + id,
			},
		},
		{
			models.InlineKeyboardButton{
				Text:         "Инструкция",
				CallbackData: "instructions",
			},
		},
		{
			models.InlineKeyboardButton{
				Text:         "На главную",
				CallbackData: "start",
			},
		},
	}}
}
