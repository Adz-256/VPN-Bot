package api

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/jackc/pgx/v4/pgxpool"
	"log/slog"
	"os"
	"os/signal"
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
		bot.WithDefaultHandler(handler),
	}
	b, err := bot.New(a.c.Token(), opts...)
	if err != nil {
		panic(err)
	}

	a.b = b

	b.Start(ctx)

	return nil
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	fmt.Println(update.Message.Chat.FirstName)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   update.Message.Text,
	})
}
