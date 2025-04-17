package wg

import (
	"fmt"
	"os"
	"text/template"
)

const wgInterfaceConfigTemplate = `[Interface]
PrivateKey = {{ .PrivateKey }}
Address = {{ .Address }}
ListenPort = {{ .ListenPort }}`

const wgPeerConfigTemplate = `[Peer]
PublicKey = {{ .PublicKey }}
AllowedIPs = {{ .AllowedIPs }}`

type WgInterfacedConfig struct {
	PrivateKey string
	Address    string
	ListenPort string
}

type WgPeerConfig struct {
	PublicKey string
	// PresharedKey string
	AllowedIPs string
}

func createWgInterface(priv string, addr string, port string) error {
	cfg := WgInterfacedConfig{
		PrivateKey: priv,
		Address:    addr,
		ListenPort: port,
	}

	tmpl, err := template.New("wgInterfaceConfig").Parse(wgInterfaceConfigTemplate)
	if err != nil {
		return fmt.Errorf("cannot create template: %v", err)
	}
	err = tmpl.Execute(os.Stdout, cfg)
	if err != nil {
		return fmt.Errorf("cannot execute template with given data: %v", err)
	}

	return nil
}

func createWgPeer(pub string, ip string) error {
	cfg := WgPeerConfig{
		PublicKey:  pub,
		AllowedIPs: ip,
	}

	tmpl, err := template.New("wgPeerConfig").Parse(wgPeerConfigTemplate)
	if err != nil {
		return fmt.Errorf("cannot create template: %v", err)
	}
	err = tmpl.Execute(os.Stdout, cfg)
	if err != nil {
		return fmt.Errorf("cannot execute template with given data: %v", err)
	}

	return nil
}
