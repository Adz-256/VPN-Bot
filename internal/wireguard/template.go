package wg

const wgInterfaceConfigTemplate = `[Interface]
PrivateKey = {{ .PrivateKey }}
Address = {{ .Address }}
ListenPort = {{ .ListenPort }}

`

const wgPeerConfigTemplate = `[Peer]
PublicKey = {{ .PublicKey }}
AllowedIPs = {{ .AllowedIPs }}

`

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
