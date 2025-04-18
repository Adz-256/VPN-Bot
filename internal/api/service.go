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

func (a *API) handleSubscriptions(ctx context.Context, b *bot.Bot, update *models.Update) {}

func (a *API) handleSupport(ctx context.Context, b *bot.Bot, update *models.Update) {
	a.l.Debug("handleBuy callback", slog.Any("chat_id", update.CallbackQuery.From.ID))

	callback := update.CallbackQuery

	_, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      callback.From.ID,
		MessageID:   callback.Message.Message.ID,
		Text:        text.Support,
		ReplyMarkup: keyboards.Support,
	})
	if err != nil {
		a.l.Error("EditMessageText error", slog.Any("error", err))
	}
}

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

// func (a *API) handleAbout(ctx context.Context, b *bot.Bot, update *models.Update) {
// 	b.SendMessage(ctx, &bot.SendMessageParams{
// 		ChatID:      update.Message.Chat.ID,
// 		Text:        text.About,
// 		ReplyMarkup: keyboards.About,
// 	})
// }

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Прости, я тебя не понял...",
	})
}
