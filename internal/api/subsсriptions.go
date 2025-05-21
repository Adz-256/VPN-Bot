package api

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/Adz-256/cheapVPN/internal/api/keyboards"
	"github.com/Adz-256/cheapVPN/internal/api/text"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	qrcode "github.com/skip2/go-qrcode"
)

func (a *API) handleSubscriptions(ctx context.Context, b *bot.Bot, update *models.Update) {
	callback := update.CallbackQuery
	slog.Debug("handleSubscriptions", slog.Any("chat_id", update.CallbackQuery.From.ID))

	accs, err := a.sub.GetUserAccounts(ctx, update.CallbackQuery.From.ID)
	if err != nil {
		slog.Error("GetUserAccounts error", slog.Any("error", err))
		return
	}

	_, err = b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      callback.From.ID,
		MessageID:   callback.Message.Message.ID,
		Text:        text.Subscriptions,
		ReplyMarkup: keyboards.Subscriptions(accs),
	})
	if err != nil {
		slog.Error("SendMessage error", slog.Any("error", err))
	}
}

func (a *API) handleShow(ctx context.Context, b *bot.Bot, update *models.Update) {
	callback := update.CallbackQuery
	slog.Debug("handleShow", slog.Any("chat_id", update.CallbackQuery.From.ID))

	_, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      callback.From.ID,
		MessageID:   callback.Message.Message.ID,
		Text:        text.Show(strings.Split(callback.Data, "_")[1]),
		ReplyMarkup: keyboards.SubsriptionConfig(strings.Split(callback.Data, "_")[1]),
	})
	if err != nil {
		slog.Error("EditMessageText error", slog.Any("error", err))
	}
}

func (a *API) handleFileRequest(ctx context.Context, b *bot.Bot, update *models.Update) {
	callback := update.CallbackQuery
	slog.Debug("handleFileRequest", slog.Any("chat_id", update.CallbackQuery.From.ID))

	accs, err := a.sub.GetUserAccounts(ctx, update.CallbackQuery.From.ID)
	if err != nil {
		slog.Error("GetUserAccounts error", slog.Any("error", err))
		return
	}
	var path string

	for _, acc := range *accs {
		if strconv.FormatInt(acc.ID, 10) == strings.Split(callback.Data, "_")[1] {
			path = acc.ConfigFile
		}
	}
	f, err := os.Open(path)
	if err != nil {
		slog.Error("Open error", slog.Any("error", err))
		return
	}
	defer f.Close()

	_, err = b.SendDocument(ctx, &bot.SendDocumentParams{
		ChatID: callback.From.ID,
		Document: &models.InputFileUpload{
			Filename: callback.Data,
			Data:     f,
		},
	})
	if err != nil {
		slog.Error("SendDocument error", slog.Any("error", err))
	}

	_, err = b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    callback.From.ID,
		MessageID: callback.Message.Message.ID,
	})
	if err != nil {
		if !strings.Contains(err.Error(), "message to delete not found") {
			slog.Error("DeleteMessage error", slog.Any("error", err))
		}
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      callback.From.ID,
		Text:        text.Start,
		ReplyMarkup: keyboards.Start,
	})
	if err != nil {
		slog.Error("SendMessage error", slog.Any("error", err))
	}
}

func (a *API) handleQRRequst(ctx context.Context, b *bot.Bot, update *models.Update) {
	callback := update.CallbackQuery
	slog.Debug("handleFileRequest", slog.Any("chat_id", update.CallbackQuery.From.ID))

	accs, err := a.sub.GetUserAccounts(ctx, update.CallbackQuery.From.ID)
	if err != nil {
		slog.Error("GetUserAccounts error", slog.Any("error", err))
		return
	}
	var path string

	for _, acc := range *accs {
		if strconv.FormatInt(acc.ID, 10) == strings.Split(callback.Data, "_")[1] {
			path = acc.ConfigFile
		}
	}
	f, err := os.Open(path)
	if err != nil {
		slog.Error("Open error", slog.Any("error", err))
		return
	}
	defer f.Close()

	body, err := io.ReadAll(f)
	if err != nil {
		slog.Error("ReadAll error", slog.Any("error", err))
		return
	}

	qrCode, err := qrcode.Encode(string(body), qrcode.Medium, 256)
	if err != nil {
		slog.Error("Encode error", slog.Any("error", err))
	}
	_, err = b.SendPhoto(ctx, &bot.SendPhotoParams{
		ChatID: callback.From.ID,
		Photo: &models.InputFileUpload{
			Filename: fmt.Sprint(callback.Data, ".conf"),
			Data:     bytes.NewReader(qrCode),
		},
	})
	if err != nil {
		slog.Error("SendDocument error", slog.Any("error", err))
	}

	_, err = b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    callback.From.ID,
		MessageID: callback.Message.Message.ID,
	})
	if err != nil {
		if !strings.Contains(err.Error(), "message to delete not found") {
			slog.Error("DeleteMessage error", slog.Any("error", err))
		}
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      callback.From.ID,
		Text:        text.Start,
		ReplyMarkup: keyboards.Start,
	})
	if err != nil {
		slog.Error("SendMessage error", slog.Any("error", err))
	}
}
