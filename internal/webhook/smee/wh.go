package smee

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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

func (w *wh) Run(ch chan map[string]any) {

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
		ch <- result
	})

	addr := fmt.Sprintf(":%s", w.port)
	log.Printf("Listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
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
