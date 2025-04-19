package api

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/Adz-256/cheapVPN/internal/api/keyboards"
	"github.com/Adz-256/cheapVPN/internal/api/text"
	"github.com/Adz-256/cheapVPN/internal/models"
	"github.com/go-telegram/bot"
	tgModels "github.com/go-telegram/bot/models"
)

func (a *API) handleBuyChooseServer(ctx context.Context, b *bot.Bot, update *tgModels.Update) {
	a.l.Debug("handleBuy callback", slog.Any("chat_id", update.CallbackQuery.From.ID))

	callback := update.CallbackQuery

	_, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      callback.From.ID,
		MessageID:   callback.Message.Message.ID,
		Text:        text.BuyChooseServer,
		ReplyMarkup: keyboards.BuyChooseServer,
	})

	if err != nil {
		a.l.Error("EditMessageText error", slog.Any("error", err))
	}
}

func (a *API) handleBuyChoosePlan(ctx context.Context, b *bot.Bot, update *tgModels.Update) {
	a.l.Debug("handleBuy callback", slog.Any("user_id", update.CallbackQuery.From.ID), slog.Any("data", update.CallbackQuery.Data))

	callback := update.CallbackQuery

	cntry := strings.Split(callback.Data, "_")[1]

	plans, err := a.s.plan.GetAllByCounty(ctx, cntry)
	if err != nil {
		a.l.Error("GetAllByCounty error", slog.Any("error", err))
	}

	kb := [][]tgModels.InlineKeyboardButton{}

	for i := 0; i < len(*plans); i++ {
		kb = append(kb, []tgModels.InlineKeyboardButton{
			{
				Text:         fmt.Sprintf("%d месяц - %s", (*plans)[i].DurationDays/32+1, (*plans)[i].Price.String()),
				CallbackData: "payment" + "_" + callback.Data + "_" + fmt.Sprintf("%d", (*plans)[i].ID),
			},
		})
	}
	kb = append(kb, []tgModels.InlineKeyboardButton{{Text: "Назад", CallbackData: "buy"}})

	_, err = b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      callback.From.ID,
		MessageID:   callback.Message.Message.ID,
		Text:        text.BuyChooseServer,
		ReplyMarkup: tgModels.InlineKeyboardMarkup{InlineKeyboard: kb},
	})

	if err != nil {
		a.l.Error("EditMessageText error", slog.Any("error", err))
	}

}

func (a *API) handlePrePayment(ctx context.Context, b *bot.Bot, update *tgModels.Update) {
	a.l.Debug("handleBuy callback", slog.Any("chat_id", update.CallbackQuery.From.ID), slog.Any("data", update.CallbackQuery.Data))

	callback := update.CallbackQuery

	planID := strings.Split(callback.Data, "_")[3]
	id, err := strconv.ParseInt(planID, 10, 64)
	if err != nil {
		a.l.Error("ParseInt error", slog.Any("error", err))
		return
	}

	plan, err := a.s.plan.GetOneByID(ctx, id)
	if err != nil {
		a.l.Error("GetOneByID error", slog.Any("error", err))
		return
	}

	u := models.User{
		UserID: callback.From.ID,
	}

	mPlan := models.Plan{
		ID:    plan.ID,
		Price: plan.Price,
	}

	link, transID, err := a.s.pay.CreatePayLink(ctx, u, mPlan, "")
	if err != nil {
		a.l.Error("CreatePayLink error", slog.Any("error", err))
		return
	}

	_, err = b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      callback.From.ID,
		MessageID:   callback.Message.Message.ID,
		Text:        fmt.Sprintf(text.BuyPrePayment, plan.Country, plan.DurationDays/32+1, plan.Price.String()),
		ReplyMarkup: keyboards.PrePayment(callback.Data, link, transID),
	})
	if err != nil {
		a.l.Error("EditMessageText error", slog.Any("error", err))
		return
	}
}

func (a *API) handlePaymentConfirm(ctx context.Context, b *bot.Bot, update *tgModels.Update) {
	a.l.Debug("handleBuy callback", slog.Any("chat_id", update.CallbackQuery.From.ID), slog.Any("data", update.CallbackQuery.Data))

	callback := update.CallbackQuery

	arrData := strings.Split(callback.Data, "_")
	transID := arrData[1]

	pay, err := a.s.pay.Get(ctx, transID)
	if err != nil {
		a.l.Error("Get error", slog.Any("error", err))
		return
	}

	pl, err := a.s.plan.GetOneByID(ctx, pay.PlanID)
	if err != nil {
		a.l.Error("GetOneByID error", slog.Any("error", err))
		return
	}

	t := time.Now().Add(time.Hour * 24 * time.Duration(pl.DurationDays))

	if pay.Status == "paid" {
		wg := &models.WgPeer{
			UserID: pay.UserID,
			EndAt:  t,
		}
		id, err := a.s.sub.CreateAccount(ctx, wg)
		if err != nil {
			a.l.Error("CreateAccount error", slog.Any("error", err))
			return
		}
		_, err = b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      callback.From.ID,
			MessageID:   callback.Message.Message.ID,
			Text:        text.BuySuccess,
			ReplyMarkup: keyboards.BuySuccess(strconv.FormatInt(id, 10)),
		})
		if err != nil {
			a.l.Error("EditMessageText error", slog.Any("error", err))
			return
		}
	} else if pay.Status == "canceled" {
		_, err = b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: callback.ID,
			Text:            text.BuyAlreadyCanceled,
			ShowAlert:       true,
		})
		if err != nil {
			a.l.Error("AnswerCallbackQuery error", slog.Any("error", err))
			return
		}
		_, err = b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      callback.From.ID,
			MessageID:   callback.Message.Message.ID,
			Text:        text.BuyChooseServer,
			ReplyMarkup: keyboards.BuyChooseServer,
		})
		if err != nil {
			a.l.Error("EditMessageText error", slog.Any("error", err))
			return
		}
	}
}
