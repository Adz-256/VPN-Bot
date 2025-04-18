package smee

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
)

type wh struct {
	addr string
	port string
}

type SmeeNotification struct {
	NotificationType string `json:"notification_type"`
	BillID           string `json:"bill_id"`
	Amount           string `json:"amount"`
	CodePro          string `json:"codepro"`
	WithdrawAmount   string `json:"withdraw_amount"`
	Unaccepted       string `json:"unaccepted"`
	Label            string `json:"label"`
	Datetime         string `json:"datetime"`
	Sender           string `json:"sender"`
	Sha1Hash         string `json:"sha1_hash"`
	OperationLabel   string `json:"operation_label"`
	OperationID      string `json:"operation_id"`
	Currency         string `json:"currency"`
}

func New(addr, port string) *wh {
	return &wh{
		addr: addr,
		port: port,
	}
}

func (w *wh) Run() chan map[string]any {
	ch := make(chan map[string]any, 1024)

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", w.addr, w.port))
	if err != nil {
		log.Panic(err)
	}

	for {
		conn, err := lis.Accept()
		reader := bufio.NewReader(conn)
		r, _ := http.ReadRequest(reader)
		if err != nil {
			log.Panic(err)
		}

		var result map[string]any
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &result)

		go func() {
			ch <- result
		}()
	}
}

func MapToNotification(m map[string]any) (SmeeNotification, error) {
	var n SmeeNotification
	bytes, err := json.Marshal(m)
	if err != nil {
		return n, err
	}
	err = json.Unmarshal(bytes, &n)
	return n, err
}
