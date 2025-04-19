package wireguard

import (
	"bufio"
	"fmt"
	"net"
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
	Pub                 string
	lastCreatedIP       string
	configOutPath       string
	*sync.Mutex
}

const (
	wgCommand      = "wg"
	wgQuickCommand = "wg-quick"
	upArg          = "up"
	downArg        = "down"
	pubkeyArg      = "pubkey"
	genArg         = "genkey"
	showArg        = "show"
	wgSetConf      = "setconf"
)

func New(interfaceName string, addr string, port string, configPath string, out string) *WgClient {
	return &WgClient{
		interfaceName:       interfaceName,
		addr:                addr,
		port:                port,
		configInterfacePath: configPath,
		lastCreatedIP:       "10.9.0.0/24",
		configOutPath:       out,
		Mutex:               &sync.Mutex{},
	}
}

func (w *WgClient) AddressWithMask() string {
	return w.addr
}

// Init применяется один раз для одного клиента
func (w *WgClient) Init() error {
	w.initFile()
	f, err := os.OpenFile(w.configInterfacePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
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
		return nil
	}
	_, err = exec.Command(wgQuickCommand, upArg, w.configInterfacePath).Output()
	if err != nil {
		return fmt.Errorf("cannot start wg interface: %v", err)
	}
	return nil
}

func (w *WgClient) Down() {
	exec.Command(wgQuickCommand, downArg, w.configInterfacePath).Run()
}

func (w *WgClient) CreateWgInterface() (pubInterface string, err error) {
	f, err := os.OpenFile(w.configInterfacePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
		Address:    w.addr,
		ListenPort: w.port,
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

func (w *WgClient) WriteUserConfig(privClient string, alowIP net.IPNet) (path string, err error) {
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
		ClientAlowedIP:   alowIP.String(),
		Endpoint:         w.addr + ":" + w.port,
	}

	w.Lock()
	defer w.Unlock()

	err = tmpl.Execute(f, cfg)
	return path, err
}

func (w *WgClient) CreateWgPeer() (allowedIP, privClient, pubClient string, err error) {
	f, err := os.OpenFile(w.configInterfacePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
	f, err := os.OpenFile(w.configInterfacePath, os.O_RDWR, 0644)
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
	err = os.Truncate(w.configInterfacePath, 0)
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

	return nil
}

func (w *WgClient) EnablePeer(pubKeyToEnable string, ip string) error {
	f, err := os.OpenFile(w.configInterfacePath, os.O_RDWR, 0644)
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
	err = os.Truncate(w.configInterfacePath, 0)
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

	return nil
}

func (w *WgClient) initFile() error {
	// Проверяем, существует ли файл
	if _, err := os.Stat(w.configInterfacePath); os.IsNotExist(err) {
		// Создаём файл с правами 0666 (rw-rw-rw-)
		file, err := os.OpenFile(w.configInterfacePath, os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		defer file.Close()
	}
	return nil
}
