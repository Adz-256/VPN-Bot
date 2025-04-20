package api

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/Adz-256/cheapVPN/internal/config/env"
	umoney "github.com/Adz-256/cheapVPN/internal/payment/uMoney"
	"github.com/Adz-256/cheapVPN/internal/repository/psql"
	"github.com/Adz-256/cheapVPN/internal/service"
	"github.com/Adz-256/cheapVPN/internal/service/payment"
	"github.com/Adz-256/cheapVPN/internal/service/plan"
	"github.com/Adz-256/cheapVPN/internal/service/subscription"
	"github.com/Adz-256/cheapVPN/internal/service/user"
	"github.com/Adz-256/cheapVPN/internal/webhook/smee"
	"github.com/Adz-256/cheapVPN/internal/wireguard"
	"github.com/go-telegram/bot"
	"github.com/jackc/pgx/v4/pgxpool"
)

type API struct {
	db         *pgxpool.Pool
	b          *bot.Bot
	l          *slog.Logger
	s          services
	cfg        *env.Config
	paymentsCh chan map[string]any
}

type services struct {
	plan service.PlanService
	pay  service.PaymentService
	sub  service.SubscriptionService
	user service.UserService
}

func New(pool *pgxpool.Pool, l *slog.Logger, cfg *env.Config, paymentsCh chan map[string]any) *API {
	return &API{
		db:         pool,
		l:          l,
		cfg:        cfg,
		paymentsCh: paymentsCh,
	}
}

// Запуск API
func (a *API) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
	}
	b, err := bot.New(a.cfg.Token(), opts...)
	if err != nil {
		panic(err)
	}

	a.b = b

	p := psql.NewPlans(a.db)
	pS := plan.New(p)
	a.s.plan = pS

	u := psql.NewUsers(a.db)
	uS := user.New(u)
	a.s.user = uS
	a.l.Info("users service created")

	pay := psql.NewPayments(a.db)
	um := umoney.New(a.cfg.PaymentAccount())
	payS := payment.New(pay, um)
	a.s.pay = payS
	a.l.Info("payment service created")

	sub := psql.NewWgPools(a.db)
	wg := wireguard.New(a.cfg.WGInterfaceName(), a.cfg.WGAddr(), a.cfg.WGPort(), a.cfg.WGPath(), a.cfg.WGOut())
	err = wg.Init()
	if err != nil {
		return err
	}
	defer wg.Down()
	subS := subscription.New(sub, wg)
	a.s.sub = subS

	a.registerHandlers()

	go a.paymentsApprover()
	b.Start(ctx)
	return nil
}

func (a *API) paymentsApprover() {
	for payment := range a.paymentsCh {
		notif, err := smee.MapToNotification(payment)
		fmt.Println(notif)
		if err != nil {
			a.l.Error("cannot map to notification", slog.Any("error", err))
			continue
		}
		if notif.Unaccepted != "true" {
			err = a.s.pay.ApprovePayment(context.Background(), notif.Label)
			if err != nil {
				a.l.Error("cannot approve payment", slog.Any("error", err))
				continue
			}
		}
	}
}

func (a *API) registerHandlers() {
	a.b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, a.handleStart)
	a.b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "start", bot.MatchTypeExact, a.handleStartCallback)

	//a.b.RegisterHandler(bot.HandlerTypeMessageText, "/about", bot.MatchTypeExact, a.handleAbout)
	// a.b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "about", bot.MatchTypeExact, a.handleAbout)

	//a.b.RegisterHandler(bot.HandlerTypeMessageText, "/support", bot.MatchTypeExact, a.handleSupport)
	a.b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "support", bot.MatchTypeExact, a.handleSupport)

	a.b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "test", bot.MatchTypeExact, a.handleTest)
	//a.b.RegisterHandler(bot.HandlerTypeMessageText, "/subscriptions", bot.MatchTypeExact, a.handleSubscriptions)
	a.b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "subscriptions", bot.MatchTypeExact, a.handleSubscriptions)
	a.b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "show_", bot.MatchTypePrefix, a.handleShow)
	a.b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "config_", bot.MatchTypePrefix, a.handleFileRequst)
	a.b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "qr_", bot.MatchTypePrefix, a.handleQRRequst)

	//a.b.RegisterHandler(bot.HandlerTypeMessageText, "/buy", bot.MatchTypeExact, a.handleBuy)
	a.b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "buy", bot.MatchTypeExact, a.handleBuyChooseServer)
	a.b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "buy_", bot.MatchTypePrefix, a.handleBuyChoosePlan)
	a.b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "payment_", bot.MatchTypePrefix, a.handlePrePayment)
	a.b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "confirm_", bot.MatchTypePrefix, a.handlePaymentConfirm)

	//a.b.RegisterHandler(bot.HandlerTypeMessageText, "/instructions", bot.MatchTypeExact, a.handleInstructions)
	a.b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "instructions", bot.MatchTypeExact, a.handleInstructions)
}
