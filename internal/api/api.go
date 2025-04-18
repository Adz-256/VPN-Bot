package api

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/jackc/pgx/v4/pgxpool"
)

type BotConfig interface {
	Token() string
}

type API struct {
	db *pgxpool.Pool
	b  *bot.Bot
	l  *slog.Logger
	c  BotConfig
}

func New(pool *pgxpool.Pool, l *slog.Logger, c BotConfig) *API {
	return &API{
		db: pool,
		l:  l,
		c:  c,
	}
}

// Запуск API
func (a *API) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
	}
	b, err := bot.New(a.c.Token(), opts...)
	if err != nil {
		panic(err)
	}

	a.b = b

	a.registerHandlers()

	b.Start(ctx)
	return nil
}

func (a *API) registerHandlers() {
	a.b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, a.handleStart)
	a.b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "start", bot.MatchTypeExact, a.handleStartCallback)

	//a.b.RegisterHandler(bot.HandlerTypeMessageText, "/about", bot.MatchTypeExact, a.handleAbout)
	// a.b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "about", bot.MatchTypeExact, a.handleAbout)

	//a.b.RegisterHandler(bot.HandlerTypeMessageText, "/support", bot.MatchTypeExact, a.handleSupport)
	a.b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "support", bot.MatchTypeExact, a.handleSupport)

	//a.b.RegisterHandler(bot.HandlerTypeMessageText, "/subscriptions", bot.MatchTypeExact, a.handleSubscriptions)
	a.b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "subscriptions", bot.MatchTypeExact, a.handleSubscriptions)

	//a.b.RegisterHandler(bot.HandlerTypeMessageText, "/buy", bot.MatchTypeExact, a.handleBuy)
	a.b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "buy", bot.MatchTypeExact, a.handleBuy)

	//a.b.RegisterHandler(bot.HandlerTypeMessageText, "/instructions", bot.MatchTypeExact, a.handleInstructions)
	a.b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "instructions", bot.MatchTypeExact, a.handleInstructions)
}
