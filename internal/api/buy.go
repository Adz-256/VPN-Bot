package api

import (
	"context"
	"log/slog"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (a *API) handleBuy(ctx context.Context, b *bot.Bot, update *models.Update) {
	a.l.Debug("handleBuy callback", slog.Any("chat_id", update.CallbackQuery.From.ID))

	callback := update.CallbackQuery

	_, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:    callback.From.ID,
		MessageID: callback.Message.Message.ID,
		Text:      "Этот бот помогает купить VPN 😊",
	})

	if err != nil {
		a.l.Error("EditMessageText error", slog.Any("error", err))
	}
}
