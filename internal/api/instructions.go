package api

import (
	"context"
	"log/slog"

	"github.com/Adz-256/cheapVPN/internal/api/keyboards"
	"github.com/Adz-256/cheapVPN/internal/api/text"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (a *API) handleInstructions(ctx context.Context, b *bot.Bot, update *models.Update) {
	a.l.Debug("handleInstructions callback", slog.Any("chat_id", update.CallbackQuery.From.ID))

	callback := update.CallbackQuery
	t := true

	_, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      callback.From.ID,
		MessageID:   callback.Message.Message.ID,
		Text:        text.Instructions,
		ReplyMarkup: keyboards.Instructions,
		ParseMode:   models.ParseModeMarkdown,
		LinkPreviewOptions: &models.LinkPreviewOptions{
			IsDisabled: &t,
		},
	})
	if err != nil {
		a.l.Error("EditMessageText error", slog.Any("error", err))
	}
}
