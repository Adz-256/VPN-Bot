package keyboards

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Adz-256/cheapVPN/internal/models"
	tgModels "github.com/go-telegram/bot/models"
)

func SubsriptionConfig(id string) tgModels.InlineKeyboardMarkup {
	return tgModels.InlineKeyboardMarkup{InlineKeyboard: [][]tgModels.InlineKeyboardButton{
		{
			tgModels.InlineKeyboardButton{
				Text:         "QR-code",
				CallbackData: "qr" + "_" + id,
			},
			tgModels.InlineKeyboardButton{
				Text:         "Конфигурационный файл",
				CallbackData: "config_" + id,
			},
		},
		{
			tgModels.InlineKeyboardButton{
				Text:         "Назад",
				CallbackData: "subscriptions",
			},
		},
	}}
}

func Subscriptions(accs *[]models.WgPeer) tgModels.InlineKeyboardMarkup {
	kb := [][]tgModels.InlineKeyboardButton{}

	for i := 0; i < len(*accs); i++ {
		var text string
		if (*accs)[i].EndAt.Unix() <= time.Now().Unix() {
			text = fmt.Sprintf("❌ %s - До %s", (*accs)[i].Name, (*accs)[i].EndAt.Format("02.01.2006"))
		} else {
			text = fmt.Sprintf("✅ %s - До %s", (*accs)[i].Name, (*accs)[i].EndAt.Format("02.01.2006"))
		}

		kb = append(kb, []tgModels.InlineKeyboardButton{
			{
				Text:         text,
				CallbackData: "show_" + strconv.FormatInt((*accs)[i].ID, 10),
			},
		})
	}
	kb = append(kb, []tgModels.InlineKeyboardButton{{Text: "Назад", CallbackData: "start"}})

	return tgModels.InlineKeyboardMarkup{InlineKeyboard: kb}
}
