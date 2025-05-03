package app

import (
	"context"
	"github.com/Adz-256/cheapVPN/internal/api"
	"github.com/Adz-256/cheapVPN/internal/closer"
	"github.com/Adz-256/cheapVPN/internal/config/env"
	"log/slog"
)

const (
	configPath = ".env"
)

type App struct {
	serviceProvider *serviceProvider

	api *api.API
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	slog.Info("starting app")
	return a.api.Run()
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initAPI,
		a.startBackgroundWorkers,
	}

	for _, f := range inits {
		if err := f(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	err := env.Load(configPath)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) startBackgroundWorkers(ctx context.Context) error {
	go a.serviceProvider.SubscriptionService(ctx).StartExpireCRON()
	go a.serviceProvider.PaymentService(ctx).StartPaymentsApprover()

	return nil
}

func (a *App) initAPI(ctx context.Context) error {
	a.api = api.New(
		a.serviceProvider.Bot(ctx),
		a.serviceProvider.PlanService(ctx),
		a.serviceProvider.PaymentService(ctx),
		a.serviceProvider.SubscriptionService(ctx),
		a.serviceProvider.UserService(ctx),
	)

	return nil
}
