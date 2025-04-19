package wireguard

const wgInterfaceConfigTemplate = `[Interface]
PrivateKey = {{ .PrivateKey }}
Address = {{ .Address }}
ListenPort = {{ .ListenPort }}

`

const wgPeerConfigTemplate = `[Peer]
PublicKey = {{ .PublicKey }}
AllowedIPs = {{ .AllowedIPs }}

`

const wgUserConfigTemplate = `[Interface]
PrivateKey = {{ .ClientPrivateKey }}
Address = {{ .ClientAlowedIP }}
DNS = 1.1.1.1

[Peer]
PublicKey = {{ .ServerPublicKey }}
AllowedIPs = 0.0.0.0/0, ::0/0
Endpoint = {{ .Endpoint }}


`

type wgInterfacedConfig struct {
	PrivateKey string
	Address    string
	ListenPort string
}

type wgPeerConfig struct {
	PublicKey string
	// PresharedKey string
	AllowedIPs string
}

type wgUserConfig struct {
	ServerPublicKey  string
	ClientPrivateKey string
	ClientAlowedIP   string
	Endpoint         string
}
