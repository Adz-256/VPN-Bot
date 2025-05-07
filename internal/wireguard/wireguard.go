package wireguard

import (
	"bufio"
	"fmt"
	"github.com/Adz-256/cheapVPN/internal/config"
	"os"
	"os/exec"
	"strings"
	"sync"
	"text/template"

	"github.com/Adz-256/cheapVPN/utils"
)

type WgClient struct {
	interfaceName       string
	configInterfacePath string
	addr                string
	port                string
	portOut             string
	Pub                 string
	lastCreatedIP       string
	configOutPath       string
	*sync.Mutex
}

const (
	wgCommand      = "wg"
	syncconfArg    = "syncconf"
	wgQuickCommand = "wg-quick"
	upArg          = "up"
	downArg        = "down"
	pubkeyArg      = "pubkey"
	genArg         = "genkey"
	showArg        = "show"
	wgSetConf      = "setconf"
	postUP         = "PostUp"
	postUpcmds     = "iptables -t nat -A POSTROUTING -s %s -o eth0 -j MASQUERADE; iptables -A INPUT -p udp -m udp --dport 51820 -j ACCEPT; iptables -A FORWARD -i %s -j ACCEPT; iptables -A FORWARD -o %s -j ACCEPT; iptables -t nat -A PREROUTING -p udp --dport 51825 -j REDIRECT --to-port 51820" //allowedips with mask; port; interfacename; intefacename
)

func New(cfg config.WgConfig) *WgClient {
	return &WgClient{
		interfaceName:       cfg.InterfaceName(),
		addr:                cfg.Address(),
		port:                cfg.Port(),
		portOut:             cfg.ExternalPort(),
		configInterfacePath: cfg.InterfaceName(),
		lastCreatedIP:       "10.9.0.1/32",
		configOutPath:       cfg.OutFilePath(),
		Mutex:               &sync.Mutex{},
	}
}

func (w *WgClient) AddressWithMask() string {
	return w.addr
}

// Init применяется один раз для одного клиента
func (w *WgClient) Init() error {
	if err := w.initFiles(); err != nil {
		return err
	}
	f, err := os.OpenFile(fmt.Sprintf("%s/%s.conf", w.configInterfacePath, w.interfaceName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		return fmt.Errorf("cannot open file: %v", err)
	}
	defer f.Close()
	s, _ := f.Stat()
	if s.Size() == 0 {
		pub, err := w.CreateWgInterface()
		if err != nil {
			return fmt.Errorf("cannot create wg interface: %v", err)
		}

		w.Pub = pub
	}
	out, err := exec.Command(wgQuickCommand, upArg, w.configInterfacePath+"/"+w.interfaceName+".conf").Output()
	if err != nil {
		return fmt.Errorf("cannot start wg interface: %v %s", out, err)
	}
	return nil
}

func (w *WgClient) Down() {
	exec.Command(wgQuickCommand, downArg, w.configInterfacePath).Run()
}

func (w *WgClient) CreateWgInterface() (pubInterface string, err error) {
	f, err := os.OpenFile(w.configInterfacePath+"/"+w.interfaceName+".conf", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		return "", fmt.Errorf("cannot open file: %v", err)
	}
	defer f.Close()

	pub, priv, err := generateKeys()
	if err != nil {
		return "", fmt.Errorf("cannot generate keys: %v", err)
	}

	cfg := wgInterfacedConfig{
		PrivateKey: priv,
		Address:    "10.9.0.1/24",
		ListenPort: w.port,
		PostUp:     fmt.Sprintf(postUpcmds, "10.9.0.0/24", w.interfaceName, w.interfaceName),
	}

	tmpl, err := template.New("wgInterfaceConfig").Parse(wgInterfaceConfigTemplate)
	if err != nil {
		return "", fmt.Errorf("cannot create template: %v", err)
	}

	w.Lock()
	defer w.Unlock()

	err = tmpl.Execute(f, cfg)
	if err != nil {
		return "", fmt.Errorf("cannot execute template with given data: %v", err)
	}

	return pub, nil
}

func (w *WgClient) WriteUserConfig(privClient string, alowIP string) (path string, err error) {
	safePriv := strings.ReplaceAll(privClient, "/", "_")
	path = w.configOutPath + "/" + safePriv + ".conf"

	// Гарантируем, что директория существует
	err = os.MkdirAll(w.configOutPath, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("cannot create config directory: %v", err)
	}

	f, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("cannot open file: %v", err)
	}
	defer f.Close()

	tmpl, err := template.New("wgUserConfig").Parse(wgUserConfigTemplate)
	if err != nil {
		return "", fmt.Errorf("cannot create template: %v", err)
	}

	cfg := wgUserConfig{
		ServerPublicKey:  w.Pub,
		ClientPrivateKey: privClient,
		ClientAllowedIP:  alowIP,
		Endpoint:         w.addr + ":" + w.portOut,
	}

	w.Lock()
	defer w.Unlock()

	err = tmpl.Execute(f, cfg)
	return path, err
}

func (w *WgClient) CreateWgPeer() (allowedIP, privClient, pubClient string, err error) {
	f, err := os.OpenFile(w.configInterfacePath+"/"+w.interfaceName+".conf", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		return "", "", "", fmt.Errorf("cannot open file: %v", err)
	}
	defer f.Close()
	pubClient, privClient, err = generateKeys()
	if err != nil {
		return "", "", "", fmt.Errorf("cannot generate keys: %v", err)
	}

	w.Lock()
	allowedIP, err = utils.IncrIP(w.lastCreatedIP)
	if err != nil {
		return "", "", "", fmt.Errorf("cannot increment ip: %v", err)
	}
	w.lastCreatedIP = allowedIP
	w.Unlock()

	cfg := wgPeerConfig{
		PublicKey:  pubClient,
		AllowedIPs: allowedIP,
	}

	tmpl, err := template.New("wgPeerConfig").Parse(wgPeerConfigTemplate)
	if err != nil {
		return "", "", "", fmt.Errorf("cannot create template: %v", err)
	}
	err = tmpl.Execute(f, cfg)
	if err != nil {
		return "", "", "", fmt.Errorf("cannot execute template with given data: %v", err)
	}
	err = w.wgSyncConfig()
	if err != nil {
		return "", "", "", fmt.Errorf("cannot sync config: %v", err)
	}

	return allowedIP, privClient, pubClient, nil
}

func generateKeys() (pub, priv string, err error) {
	privateKeyBytes, err := exec.Command(wgCommand, genArg).Output()
	if err != nil {
		return "", "", fmt.Errorf("cannot generate private key: %v", err)
	}
	privateKey := strings.TrimSpace(string(privateKeyBytes))

	// Генерация публичного ключа из приватного
	cmd := exec.Command(wgCommand, pubkeyArg)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", "", fmt.Errorf("cannot create stdin pipe: %v", err)
	}
	go func() {
		defer stdin.Close()
		stdin.Write([]byte(privateKey))
	}()
	publicKeyBytes, err := cmd.Output()
	if err != nil {
		return "", "", fmt.Errorf("cannot generate public key: %v", err)
	}
	publicKey := strings.TrimSpace(string(publicKeyBytes))

	return publicKey, privateKey, nil
}

// RemovePeer удаляет пир из конфигурационного файла wg0.conf по публичному ключу.
func (w *WgClient) BlockPeer(pubKeyToRemove string) error {
	f, err := os.OpenFile(w.configInterfacePath+"/"+w.interfaceName+".conf", os.O_RDWR, 0777)
	if err != nil {
		return fmt.Errorf("cannot open file: %v", err)
	}
	defer f.Close()

	found := false
	// Читаем файл построчно
	w.Lock()
	defer w.Unlock()
	scanner := bufio.NewScanner(f)
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
		if line == fmt.Sprintf("PublicKey = %s", pubKeyToRemove) {
			found = true
		}

		if found && strings.Contains(line, "AllowedIPs") {
			arr := strings.Split(line, " ")
			arr[2] = "0.0.0.0/0"
			lines[len(lines)-1] = strings.Join(arr, " ")
			found = false
		}
	}

	// Проверяем на ошибки чтения
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading the file: %v", err)
	}

	// Перезаписываем файл с удалённым пиром
	err = os.Truncate(w.configInterfacePath+"/"+w.interfaceName+".conf", 0)
	if err != nil {
		return fmt.Errorf("cannot truncate file: %v", err)
	}

	f.Seek(0, 0)

	// Записываем обновлённые данные обратно в файл
	writer := bufio.NewWriter(f)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return fmt.Errorf("error writing to file: %v", err)
		}
	}

	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("error flushing buffer: %v", err)
	}

	err = w.wgSyncConfig()
	if err != nil {
		return fmt.Errorf("cannot sync config: %v", err)
	}

	return nil
}

func (w *WgClient) EnablePeer(pubKeyToEnable string, ip string) error {
	f, err := os.OpenFile(w.configInterfacePath+"/"+w.interfaceName+".conf", os.O_RDWR, 0777)
	if err != nil {
		return fmt.Errorf("cannot open file: %v", err)
	}
	defer f.Close()

	found := false
	// Читаем файл построчно
	scanner := bufio.NewScanner(f)
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
		if line == fmt.Sprintf("PublicKey = %s", pubKeyToEnable) {
			found = true
		}

		if found && strings.Contains(line, "AllowedIPs") {
			arr := strings.Split(line, " ")
			arr[2] = ip
			lines[len(lines)-1] = strings.Join(arr, " ")
			found = false
		}
	}

	// Проверяем на ошибки чтения
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading the file: %v", err)
	}

	// Перезаписываем файл с удалённым пиром
	err = os.Truncate(w.configInterfacePath+"/"+w.interfaceName+".conf", 0)
	if err != nil {
		return fmt.Errorf("cannot truncate file: %v", err)
	}

	f.Seek(0, 0)

	// Записываем обновлённые данные обратно в файл
	writer := bufio.NewWriter(f)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return fmt.Errorf("error writing to file: %v", err)
		}
	}

	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("error flushing buffer: %v", err)
	}
	err = w.wgSyncConfig()
	if err != nil {
		return fmt.Errorf("cannot sync config: %v", err)
	}
	return nil
}

func (w *WgClient) initFiles() error {
	// Проверяем, существует ли директория
	if err := os.MkdirAll(w.configInterfacePath, 0777); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Полный путь до файла
	filePath := w.configInterfacePath + "/" + w.interfaceName + ".conf"

	// Проверяем, существует ли файл
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Создаём файл с правами 0666 (rw-rw-rw-)
		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0777)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		defer file.Close()
	}

	// Создаём файл конфига для пользователей
	if err := os.MkdirAll(w.configOutPath, 0777); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	return nil
}

func (w *WgClient) wgSyncConfig() error {
	out, err := exec.Command("wg-quick", "strip", w.configInterfacePath+"/"+w.interfaceName+".conf").Output()
	if err != nil {
		return err
	}
	f, _ := os.OpenFile("config/temp.conf", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0777)
	defer f.Close()

	f.Write(out)
	return exec.Command("wg", "syncconf", "wg1", "config/temp.conf").Run()
} //wg syncconf wg1 <(wg-quick strip wg1)
