package smee

import (
	"encoding/json"
	"fmt"
	"github.com/Adz-256/cheapVPN/internal/config"
	"io"
	"log"
	"net/http"
)

type WH struct {
	addr string
	port string
	ch   chan map[string]any
}

type Notification struct {
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

func New(cfg config.WhConfig) *WH {
	return &WH{
		addr: cfg.Address(),
		port: cfg.Port(),
		ch:   make(chan map[string]any, 1024),
	}
}

func (w *WH) Run() {

	http.HandleFunc("/", func(wr http.ResponseWriter, r *http.Request) {
		var result map[string]any
		defer r.Body.Close()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("error reading body: %v", err)
			return
		}
		err = json.Unmarshal(body, &result)
		if err != nil {
			log.Printf("error unmarshaling json: %v", err)
			return
		}
		log.Printf("Received webhook: %+v", result)
		w.ch <- result
	})

	addr := fmt.Sprintf(":%s", w.port)
	log.Printf("Listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func (w *WH) Channel() <-chan map[string]any {
	return w.ch
}

func MapToNotification(m map[string]any) (Notification, error) {
	var n Notification
	bytes, err := json.Marshal(m)
	if err != nil {
		return n, err
	}
	err = json.Unmarshal(bytes, &n)
	return n, err
}
