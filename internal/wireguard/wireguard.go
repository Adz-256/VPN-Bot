package wg

import (
	"os"
	"os/exec"
	"strings"
)

type WgClinet struct {
	ip   string
	port string
	f    *os.File
}

const (
	wgCommand = "wg"
	pubkeyArg = "pubkey"
	genArg    = "gen"
	showArg   = "show"
)

func New(ip string, port string, configPath string) *WgClinet {
	out, _ := exec.Command(wgCommand, showArg).Output()
	if len(out) == 0 {
		//initWgInterface()
	}
	// f, err := os.Open(path)
	// if err != nil {
	// 	panic(err)
	// }
	return &WgClinet{
		ip:   ip,
		port: port,
		// f:    f,
	}
}

/*func initWgInterface() error {
	pub, priv := generateKeys()
	f, err := os.Create("wg0.conf")
	if err != nil {
		return fmt.Errorf("cannot create file: %v", err)
	}

	f.WriteString()
}*/

func generateKeys() (pub, priv string) {
	privateKeyBytes, err := exec.Command(wgCommand, genArg).Output()
	if err != nil {
		panic(err)
	}
	privateKey := strings.TrimSpace(string(privateKeyBytes))

	// Генерация публичного ключа из приватного
	cmd := exec.Command(wgCommand, pubkeyArg)
	stdin, _ := cmd.StdinPipe()
	go func() {
		defer stdin.Close()
		stdin.Write([]byte(privateKey))
	}()
	publicKeyBytes, _ := cmd.Output()
	publicKey := strings.TrimSpace(string(publicKeyBytes))
	return publicKey, privateKey
}
