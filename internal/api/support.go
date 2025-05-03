package api

import (
	"context"
	"log/slog"

	"github.com/Adz-256/cheapVPN/internal/api/keyboards"
	"github.com/Adz-256/cheapVPN/internal/api/text"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (a *API) handleSupport(ctx context.Context, b *bot.Bot, update *models.Update) {
	slog.Debug("handleBuy callback", slog.Any("chat_id", update.CallbackQuery.From.ID))

	callback := update.CallbackQuery

	_, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      callback.From.ID,
		MessageID:   callback.Message.Message.ID,
		Text:        text.Support,
		ReplyMarkup: keyboards.Support,
	})
	if err != nil {
		slog.Error("EditMessageText error", slog.Any("error", err))
	}
}
