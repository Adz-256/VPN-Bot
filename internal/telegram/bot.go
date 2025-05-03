package telegram

import (
	"context"
	"github.com/Adz-256/cheapVPN/internal/config"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func InitBotWithDefaultHandler(cfg config.BotConfig) *bot.Bot {
	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
	}
	b, err := bot.New(cfg.Token(), opts...)
	if err != nil {
		panic(err)
	}

	return b
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Прости, я тебя не понял...",
	})
}
