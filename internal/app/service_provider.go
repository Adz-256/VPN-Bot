package app

import (
	"context"
	"github.com/Adz-256/cheapVPN/internal/broker"
	"github.com/Adz-256/cheapVPN/internal/broker/kafka"
	"github.com/Adz-256/cheapVPN/internal/metrics"
	"github.com/Adz-256/cheapVPN/internal/metrics/prometheus"

	"github.com/Adz-256/cheapVPN/internal/api"
	"github.com/Adz-256/cheapVPN/internal/closer"
	"github.com/Adz-256/cheapVPN/internal/config"
	"github.com/Adz-256/cheapVPN/internal/config/env"
	proc "github.com/Adz-256/cheapVPN/internal/payment"
	umoney "github.com/Adz-256/cheapVPN/internal/payment/uMoney"
	"github.com/Adz-256/cheapVPN/internal/repository"
	"github.com/Adz-256/cheapVPN/internal/repository/psql"
	"github.com/Adz-256/cheapVPN/internal/service"
	"github.com/Adz-256/cheapVPN/internal/service/payment"
	"github.com/Adz-256/cheapVPN/internal/service/plan"
	"github.com/Adz-256/cheapVPN/internal/service/subscription"
	"github.com/Adz-256/cheapVPN/internal/service/user"
	"github.com/Adz-256/cheapVPN/internal/telegram"
	"github.com/Adz-256/cheapVPN/internal/webhook"
	"github.com/Adz-256/cheapVPN/internal/webhook/smee"
	"github.com/Adz-256/cheapVPN/internal/wireguard"
	"github.com/Adz-256/cheapVPN/pkg/clients/postgres"
	"github.com/go-telegram/bot"
	"github.com/jackc/pgx/v4/pgxpool"
)

type serviceProvider struct {
	logCfg    config.LoggerConfig
	dbCfg     config.DBConfig
	payCfg    config.PaymentConfig
	wgCfg     config.WgConfig
	whCfg     config.WhConfig
	botCfg    config.BotConfig
	subCfg    config.SubscriptionConfig
	brokerCfg config.BrokerConfig

	db  *pgxpool.Pool
	bot *bot.Bot

	paymentBroker broker.Broker

	paymentsRepo repository.PaymentRepository

	userRepo repository.UserRepository

	wgClient   *wireguard.WgClient
	wgPoolRepo repository.WgPoolRepository

	planRepo repository.PlanRepository

	paymentConsumer  broker.Consumer
	paymentPublisher broker.Publisher
	paymentWebhook   webhook.Webhook
	paymentProcessor proc.Payment
	paymentService   service.PaymentService

	userService service.UserService

	subService service.SubscriptionService

	planService service.PlanService

	api *api.API

	metricsServer metrics.Server
}

const (
	paymentTopic   = "payments"
	paymentGroupID = "payment_group"
)

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) BrokerConfig() config.BrokerConfig {
	if s.brokerCfg == nil {
		brokerCfg, err := env.NewPaymentBrokerConfig()
		if err != nil {
			panic(err)
		}
		s.brokerCfg = brokerCfg
	}

	return s.brokerCfg
}

func (s *serviceProvider) DBConfig() config.DBConfig {
	if s.dbCfg == nil {
		dbCfg, err := env.NewPGConfig()
		if err != nil {
			panic(err)
		}
		s.dbCfg = dbCfg
	}

	return s.dbCfg
}

func (s *serviceProvider) SubConfig() config.SubscriptionConfig {
	if s.subCfg == nil {
		s.subCfg = env.NewSubscription()
	}

	return s.subCfg
}

func (s *serviceProvider) PayConfig() config.PaymentConfig {
	if s.payCfg == nil {
		payCfg, err := env.NewPaymentConfig()
		if err != nil {
			panic(err)
		}
		s.payCfg = payCfg
	}
	return s.payCfg
}

func (s *serviceProvider) WgConfig() config.WgConfig {
	if s.wgCfg == nil {
		wgCfg, err := env.NewWGConfig()
		if err != nil {
			panic(err)
		}
		s.wgCfg = wgCfg
	}

	return s.wgCfg
}

func (s *serviceProvider) LoggerConfig() config.LoggerConfig {
	if s.logCfg == nil {
		logCfg, err := env.NewLoggerConfig()
		if err != nil {
			panic(err)
		}
		s.logCfg = logCfg
	}

	return s.logCfg
}

func (s *serviceProvider) WHConfig() config.WhConfig {
	if s.whCfg == nil {
		whCfg, err := env.NewWebhookConfig()
		if err != nil {
			panic(err)
		}
		s.whCfg = whCfg
	}

	return s.whCfg
}

func (s *serviceProvider) BotConfig() config.BotConfig {
	if s.botCfg == nil {
		botCfg, err := env.NewBotConfig()
		if err != nil {
			panic(err)
		}
		s.botCfg = botCfg
	}

	return s.botCfg
}

func (s *serviceProvider) DBClient(ctx context.Context) *pgxpool.Pool {
	if s.dbCfg == nil {
		pool := postgres.New(ctx, s.DBConfig())
		defer closer.Add(func() error {
			pool.Close()
			return nil
		},
		)

		err := pool.Ping(ctx)
		if err != nil {
			panic(err)
		}

		s.db = pool
	}

	return s.db
}

func (s *serviceProvider) MetricsServer(_ context.Context) metrics.Server {
	if s.metricsServer == nil {
		s.metricsServer = prometheus.New()
	}

	return s.metricsServer
}

func (s *serviceProvider) BrokerClient(_ context.Context) broker.Broker {
	if s.paymentBroker == nil {
		b, err := kafka.New(s.BrokerConfig())
		if err != nil {
			panic(err)
		}

		s.paymentBroker = b
	}

	return s.paymentBroker
}

func (s *serviceProvider) PaymentConsumer(ctx context.Context) broker.Consumer {
	if s.paymentConsumer == nil {
		s.paymentConsumer = s.BrokerClient(ctx).NewReader(paymentGroupID, paymentTopic)
	}

	return s.paymentConsumer
}

func (s *serviceProvider) PaymentPublisher(ctx context.Context) broker.Publisher {
	if s.paymentPublisher == nil {
		s.paymentPublisher = s.BrokerClient(ctx).NewWriter(paymentTopic)
	}

	return s.paymentPublisher
}

func (s *serviceProvider) UserRepo(ctx context.Context) repository.UserRepository {
	if s.userRepo == nil {
		s.userRepo = psql.NewUsers(s.DBClient(ctx))
	}

	return s.userRepo
}

func (s *serviceProvider) PlanRepo(ctx context.Context) repository.PlanRepository {
	if s.planRepo == nil {
		s.planRepo = psql.NewPlans(s.DBClient(ctx))
	}
	return s.planRepo
}

func (s *serviceProvider) PaymentRepo(ctx context.Context) repository.PaymentRepository {
	if s.paymentsRepo == nil {
		s.paymentsRepo = psql.NewPayments(s.DBClient(ctx))
	}

	return s.paymentsRepo
}

func (s *serviceProvider) WgPoolRepo(ctx context.Context) repository.WgPoolRepository {
	if s.wgPoolRepo == nil {
		s.wgPoolRepo = psql.NewWgPools(s.DBClient(ctx))
	}

	return s.wgPoolRepo
}

func (s *serviceProvider) PlanService(ctx context.Context) service.PlanService {
	if s.planService == nil {
		s.planService = plan.NewService(s.PlanRepo(ctx))
	}
	return s.planService
}

func (s *serviceProvider) WgClient(_ context.Context) *wireguard.WgClient {
	if s.wgClient == nil {
		s.wgClient = wireguard.New(s.WgConfig())

		err := s.wgClient.Init()
		if err != nil {
			panic(err)
		}
	}

	return s.wgClient
}

func (s *serviceProvider) SubscriptionService(ctx context.Context) service.SubscriptionService {
	if s.subService == nil {
		s.subService = subscription.NewService(s.WgPoolRepo(ctx), s.WgClient(ctx), s.SubConfig())
	}

	return s.subService
}

func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = user.NewService(s.UserRepo(ctx))
	}

	return s.userService
}

func (s *serviceProvider) PaymentProcessor(_ context.Context) proc.Payment {
	if s.paymentProcessor == nil {
		s.paymentProcessor = umoney.New(s.PayConfig())
	}

	return s.paymentProcessor
}

func (s *serviceProvider) PaymentService(ctx context.Context) service.PaymentService {
	if s.paymentService == nil {
		s.paymentService = payment.NewService(s.PaymentRepo(ctx),
			s.Webhook(ctx), s.PaymentProcessor(ctx), s.PaymentConsumer(ctx))
	}

	return s.paymentService
}

func (s *serviceProvider) Bot(_ context.Context) *bot.Bot {
	if s.bot == nil {
		s.bot = telegram.InitBotWithDefaultHandler(s.BotConfig())
	}

	return s.bot
}

func (s *serviceProvider) Webhook(ctx context.Context) webhook.Webhook {
	if s.paymentWebhook == nil {
		s.paymentWebhook = smee.New(s.WHConfig(), s.PaymentPublisher(ctx))
	}

	return s.paymentWebhook
}

func (s *serviceProvider) API(ctx context.Context) *api.API {
	if s.api == nil {
		s.api = api.New(
			s.Bot(ctx),
			s.PlanService(ctx),
			s.PaymentService(ctx),
			s.SubscriptionService(ctx),
			s.UserService(ctx),
		)
	}

	return s.api
}
