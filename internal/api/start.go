package api

import (
	"context"
	"log/slog"

	"github.com/Adz-256/cheapVPN/internal/api/keyboards"
	"github.com/Adz-256/cheapVPN/internal/api/text"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (a *API) handleStart(ctx context.Context, b *bot.Bot, update *models.Update) {
	a.l.Debug("handleStart", slog.Any("chat_id", update.Message.From.ID))

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        text.Start,
		ReplyMarkup: keyboards.Start,
	})
	if err != nil {
		a.l.Error("SendMessage error", slog.Any("error", err))
	}
}

func (a *API) handleStartCallback(ctx context.Context, b *bot.Bot, update *models.Update) {
	callback := update.CallbackQuery
	a.l.Debug("handleStartCallback", slog.Any("chat_id", callback.From.ID))

	_, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      callback.From.ID,
		MessageID:   callback.Message.Message.ID,
		Text:        text.Start,
		ReplyMarkup: keyboards.Start,
	})
	if err != nil {
		a.l.Error("EditMessageText error", slog.Any("error", err))
	}
}
