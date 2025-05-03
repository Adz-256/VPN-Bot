package env

import (
	"errors"
	"github.com/Adz-256/cheapVPN/internal/config"
	"os"
)

type wgConfig struct {
	path          string
	addr          string
	port          string
	externalPort  string
	interfaceName string
	out           string
}

const (
	wgPathEnv          = "WIREGUARD_CONFIG_PATH"
	wgAddressEnv       = "WIREGUARD_ADDRESS"
	wgPortEnv          = "WIREGUARD_PORT"
	wgExternalPortEnv  = "WIREGUARD_EXTERNAL_PORT"
	wgInterfaceNameEnv = "WIREGUARD_INTERFACE_NAME"
	wgOutFilePathEnv   = "WIREGUARD_OUT"
)

var (
	ErrNoWgPath          = errors.New("no WIREGUARD_CONFIG_PATH environment variable")
	ErrNoWgAddr          = errors.New("no WIREGUARD_ADDRESS environment variable")
	ErrNoWgPort          = errors.New("no WIREGUARD_PORT environment variable")
	ErrNoWgExternalPort  = errors.New("no WIREGUARD_EXTERNAL_PORT environment variable")
	ErrNoWgInterfaceName = errors.New("no WIREGUARD_INTERFACE_NAME environment variable")
	ErrNoWgOut           = errors.New("no WIREGUARD_OUT environment variable")
)

func NewWGConfig() (config.WgConfig, error) {
	path := os.Getenv(wgPathEnv)
	if path == "" {
		return nil, ErrNoWgPath
	}
	addr := os.Getenv(wgAddressEnv)
	if addr == "" {
		return nil, ErrNoWgAddr
	}
	port := os.Getenv(wgPortEnv)
	if port == "" {
		return nil, ErrNoWgPort
	}
	externalPort := os.Getenv(wgExternalPortEnv)
	if externalPort == "" {
		return nil, ErrNoWgExternalPort
	}
	interfaceName := os.Getenv(wgInterfaceNameEnv)
	if interfaceName == "" {
		return nil, ErrNoWgInterfaceName
	}
	outPath := os.Getenv(wgOutFilePathEnv)
	if outPath == "" {
		return nil, ErrNoWgOut
	}

	return &wgConfig{
		path:          path,
		addr:          addr,
		port:          port,
		externalPort:  externalPort,
		interfaceName: interfaceName,
		out:           outPath,
	}, nil
}

func (cfg *wgConfig) Path() string {
	return cfg.path
}

func (cfg *wgConfig) Address() string {
	return cfg.addr
}

func (cfg *wgConfig) Port() string {
	return cfg.port
}

func (cfg *wgConfig) ExternalPort() string {
	return cfg.externalPort
}

func (cfg *wgConfig) InterfaceName() string {
	return cfg.interfaceName
}

func (cfg *wgConfig) OutFilePath() string {
	return cfg.out
}
