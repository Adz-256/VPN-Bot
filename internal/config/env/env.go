package env

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
)

var (
	NoTokenError = (errors.New("no token provided"))
	NoDsnError   = (errors.New("no dsn provided"))
	NoEnvError   = (errors.New("no env provided, env set as development"))
)

type Config struct {
	bot botConfig
	dsn string
	env string
}

type botConfig struct {
	token string
}

func New() (cfg *Config, err error) {
	cfg = &Config{}
	err = godotenv.Load(".env")
	if err != nil {
		return &Config{}, err
	}

	err = cfg.loadDSNConfig()
	if err != nil {
		return &Config{}, err
	}
	err = cfg.loadEnvConfig()
	if err != nil {
		return &Config{}, err
	}
	err = cfg.loadBotConfig()
	if err != nil {
		return &Config{}, err
	}

	return cfg, nil
}

func (c *Config) DSN() string {
	return c.dsn
}

func (c *Config) ENV() string {
	return c.env
}
func (c *Config) Token() string {
	return c.bot.token
}

func (c *Config) loadBotConfig() error {
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		return NoTokenError
	}
	c.bot.token = token

	return nil
}

func (c *Config) loadDSNConfig() error {
	dsn := os.Getenv("BOT_TOKEN")
	if dsn == "" {
		return NoDsnError
	}
	c.dsn = dsn

	return nil
}

func (c *Config) loadEnvConfig() error {
	env := os.Getenv("ENV")
	if env == "" {
		c.env = "development"
		return NoEnvError
	}
	c.env = env

	return nil
}
