package config

type LoggerConfig interface {
	Level() string
}

type PaymentConfig interface {
	AccountID() string
}

type WhConfig interface {
	Address() string
	Port() string
}

type BotConfig interface {
	Token() string
}

type WgConfig interface {
	Path() string
	Address() string
	Port() string
	ExternalPort() string
	InterfaceName() string
	OutFilePath() string
}

type DBConfig interface {
	DSN() string
}
