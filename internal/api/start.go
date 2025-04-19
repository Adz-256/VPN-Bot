package api

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/Adz-256/cheapVPN/internal/api/keyboards"
	"github.com/Adz-256/cheapVPN/internal/api/text"
	"github.com/Adz-256/cheapVPN/internal/models"
	"github.com/go-telegram/bot"
	tgModels "github.com/go-telegram/bot/models"
)

func (a *API) handleStart(ctx context.Context, b *bot.Bot, update *tgModels.Update) {
	a.l.Debug("handleStart", slog.Any("chat_id", update.Message.From.ID))

	u := &models.User{
		UserID:   update.Message.From.ID,
		Username: update.Message.From.FirstName,
	}

	_, err := a.s.user.Create(ctx, u)

	if err != nil && err != sql.ErrNoRows {
		a.l.Error("Create error", slog.Any("error", err))
	}
	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        text.Start,
		ReplyMarkup: keyboards.Start,
	})
	if err != nil {
		a.l.Error("SendMessage error", slog.Any("error", err))
		return
	}
}

func (a *API) handleStartCallback(ctx context.Context, b *bot.Bot, update *tgModels.Update) {
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
		return
	}
}
