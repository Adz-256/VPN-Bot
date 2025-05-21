package api

import (
	"context"
	"os"
	"os/signal"

	"github.com/Adz-256/cheapVPN/internal/service"
	"github.com/go-telegram/bot"
)

type API struct {
	b    *bot.Bot
	plan service.PlanService
	pay  service.PaymentService
	sub  service.SubscriptionService
	user service.UserService
}

func New(
	b *bot.Bot,
	plan service.PlanService,
	pay service.PaymentService,
	sub service.SubscriptionService,
	user service.UserService,
) *API {
	return &API{
		b:    b,
		plan: plan,
		pay:  pay,
		sub:  sub,
		user: user,
	}
}

// Run Запуск API
func (a *API) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	a.registerHandlers()

	a.b.Start(ctx)
	return nil
}

func (a *API) registerHandlers() {
	a.b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, a.handleStart)
	a.b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "start", bot.MatchTypeExact, a.handleStartCallback)
	a.b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "support", bot.MatchTypeExact, a.handleSupport)

	a.b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "test", bot.MatchTypeExact, a.handleTest)
	a.b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "subscriptions", bot.MatchTypeExact, a.handleSubscriptions)
	a.b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "show_", bot.MatchTypePrefix, a.handleShow)
	a.b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "config_", bot.MatchTypePrefix, a.handleFileRequest)
	a.b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "qr_", bot.MatchTypePrefix, a.handleQRRequst)

	a.b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "buy", bot.MatchTypeExact, a.handleBuyChooseServer)
	a.b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "buy_", bot.MatchTypePrefix, a.handleBuyChoosePlan)
	a.b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "payment_", bot.MatchTypePrefix, a.handlePrePayment)
	a.b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "confirm_", bot.MatchTypePrefix, a.handlePaymentConfirm)

	a.b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "instructions", bot.MatchTypeExact, a.handleInstructions)
}
