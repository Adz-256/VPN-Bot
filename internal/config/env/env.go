package env

import (
	"github.com/joho/godotenv"
)

//var (
//	ErrorNoToken          = (errors.New("no token provided"))
//	ErrorNoDsn            = (errors.New("no dsn provided"))
//	ErrorNoEnv            = (errors.New("no env provided, env set as development"))
//	ErrorWgPathIsEmpty    = (errors.New("wireguard config path is empty"))
//	ErrorNoWgPort         = (errors.New("no wireguard port provided"))
//	ErrNoWgAddress        = (errors.New("no wireguard address provided"))
//	ErrNoPaymentAcc       = (errors.New("no payment account provided"))
//	ErrorNoWgExternalPort = (errors.New("no wireguard external port provided"))
//)

//type Config struct {
//	bot botConfig
//	dsn string
//	acc string
//	env string
//	wh  webhookConfig
//	wg  wgConfig
//}

//type botConfig struct {
//	token string
//}

func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}

//
//func New() (cfg *Config, err error) {
//	cfg = &Config{}
//	err = godotenv.Load(".env")
//	if err != nil {
//		slog.Error("Error loading .env file", slog.Any("error", err))
//	}
//
//	err = cfg.loadDSNConfig()
//	if err != nil {
//		return nil, err
//	}
//	err = cfg.loadEnvConfig()
//	if err != nil {
//		return nil, err
//	}
//	err = cfg.loadBotConfig()
//	if err != nil {
//		return nil, err
//	}
//
//	err = cfg.loadWgConfig()
//	if err != nil {
//		return nil, err
//	}
//
//	cfg.loadPaymentAccount()
//
//	cfg.loadWhConfig()
//
//	return cfg, nil
//}
//
//func (c *Config) DSN() string {
//	return c.dsn
//}
//
//func (c *Config) WebhookAddress() string {
//	return c.wh.addr
//}
//
//func (c *Config) WebhookPort() string {
//	return c.wh.port
//}
//
//func (c *Config) ENV() string {
//	return c.env
//}
//func (c *Config) Token() string {
//	return c.bot.token
//}
//
//func (c *Config) WGPath() string {
//	return c.wg.path
//}
//
//func (c *Config) WGAddr() string {
//	return c.wg.addr
//}
//
//func (c *Config) WGPort() string {
//	return c.wg.port
//}
//
//func (c *Config) WGInterfaceName() string {
//	return c.wg.interfaceName
//}
//
//func (c *Config) WGOut() string {
//	return c.wg.out
//}
//
//func (c *Config) PaymentAccount() string {
//	return c.acc
//}
//
//func (c *Config) loadBotConfig() error {
//	token := os.Getenv("BOT_TOKEN")
//	if token == "" {
//		return ErrorNoToken
//	}
//	c.bot.token = token
//
//	return nil
//}
//
//func (c *Config) loadPaymentAccount() error {
//	acc := os.Getenv("PAYMENT_ACCOUNT")
//	if acc == "" {
//		return ErrNoPaymentAcc
//	}
//
//	c.acc = acc
//
//	return nil
//}
//
//func (c *Config) loadDSNConfig() error {
//	dsn := os.Getenv("DSN")
//	if dsn == "" {
//		return ErrorNoDsn
//	}
//	c.dsn = dsn
//
//	return nil
//}
//
//func (c *Config) loadEnvConfig() error {
//	env := os.Getenv("ENV")
//	if env == "" {
//		c.env = "development"
//		return ErrorNoEnv
//	}
//	c.env = env
//
//	return nil
//}
//
//func (c *Config) loadWhConfig() error {
//	addr := os.Getenv("WEBHOOK_ADDRESS")
//	if addr == "" {
//		return ErrNoWgAddress
//	}
//
//	c.wh.addr = addr
//
//	port := os.Getenv("WEBHOOK_PORT")
//	if port == "" {
//		return ErrorNoWgPort
//	}
//
//	c.wh.port = port
//
//	return nil
//}
//
//// TODO: Убрать паники и return переделать под логгер
//func (c *Config) loadWgConfig() error {
//	path := os.Getenv("WIREGUARD_CONFIG_PATH")
//	if path == "" {
//		return ErrorWgPathIsEmpty
//	}
//	c.wg.path = path
//
//	addr := os.Getenv("WIREGUARD_ADDRESS")
//	if addr == "" {
//		return ErrNoWgAddress
//	}
//	c.wg.addr = addr
//
//	port := os.Getenv("WIREGUARD_PORT")
//	if port == "" {
//		return ErrorNoWgPort
//	}
//	c.wg.port = port
//
//	externalPort := os.Getenv("WIREGUARD_EXTERNAL_PORT")
//	if externalPort == "" {
//		return ErrorNoWgExternalPort
//	}
//	c.wg.externalPort = externalPort
//
//	interfaceName := os.Getenv("WIREGUARD_INTERFACE_NAME")
//	if interfaceName == "" {
//		interfaceName = "wg0"
//	}
//	c.wg.interfaceName = interfaceName
//
//	out := os.Getenv("WIREGUARD_OUT")
//	if out == "" {
//		out = "config"
//	}
//	c.wg.out = out
//
//	return nil
//}
