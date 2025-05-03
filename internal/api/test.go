package api

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"github.com/Adz-256/cheapVPN/internal/api/keyboards"
	"github.com/Adz-256/cheapVPN/internal/api/text"
	"github.com/Adz-256/cheapVPN/internal/models"
	"github.com/go-telegram/bot"
	tgModels "github.com/go-telegram/bot/models"
)

func (a *API) handleTest(ctx context.Context, b *bot.Bot, update *tgModels.Update) {
	slog.Debug("handleTest", slog.Any("chat_id", update.CallbackQuery.From.ID))

	callback := update.CallbackQuery

	t := time.Now().Add(time.Hour * 24 * time.Duration(1))

	wg := &models.WgPeer{
		UserID: update.CallbackQuery.From.ID,
		EndAt:  t,
	}

	accs, err := a.sub.GetUserAccounts(ctx, update.CallbackQuery.From.ID)

	if err != nil || len(*accs) != 0 {
		slog.Debug("GetUserAccounts error", slog.Any("error", err), slog.Any("accs", accs))
		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: callback.ID,
			Text:            text.TestFail,
		})
		return
	}

	slog.Debug("CreateAccount", slog.Any("wg", wg))
	id, err := a.sub.CreateAccount(ctx, wg)
	if err != nil {
		slog.Error("CreateAccount error", slog.Any("error", err))
		return
	}
	_, err = b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      callback.From.ID,
		MessageID:   callback.Message.Message.ID,
		Text:        text.TestSuccess,
		ReplyMarkup: keyboards.BuySuccess(strconv.FormatInt(id, 10)),
	})
	if err != nil {
		slog.Error("EditMessageText error", slog.Any("error", err))
		return
	}
}
