package umoney

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	baseUrl       = "https://yoomoney.ru/quickpay/confirm.xml"
	receiver      = "receiver"
	quickpayForm  = "quickpay-form"
	targets       = "targets"
	paymentType   = "paymentType"
	sum           = "sum"
	transactionID = "label"
)

type Payment struct {
}

type Quickpay struct {
	Recieiver    string
	QuickpayForm string
	Targets      string
	Label        string
	PaymentType  string
	Sum          string
}

func New(umoneyAccount string) *Payment {
	return &Payment{}
}

// https://yoomoney.ru/quickpay/confirm.xml?receiver=4100117034899495&quickpay-form=shop&targets=Sponsor%20this%20project&paymentType=SB&sum=5
func (p *Payment) CreatePayLink(qp Quickpay) (link string, transID string, err error) {
	transID = qp.Label

	if transID == "" {
		transID = generateTransactionId()
	}

	payURL := configureRequestURL(qp, transID)

	// Создаем POST-запрос по ссылке
	req, err := http.NewRequest(http.MethodPost, payURL, nil)
	if err != nil {
		// логируем, если нужно
		return "", "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// логируем, если нужно
		return "", "", err
	}
	defer resp.Body.Close()

	// Example response:
	// Found. Redirecting to https://...
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		// логируем, если нужно
		return "", "", err
	}

	bodySplit := strings.Split(string(b), " ")

	url := bodySplit[4]
	return url, transID, nil
}

func configureRequestURL(qp Quickpay, transID string) string {
	values := url.Values{}
	values.Set(receiver, qp.Recieiver)
	values.Set(quickpayForm, qp.QuickpayForm)
	values.Set(targets, qp.Targets)
	values.Set(paymentType, qp.PaymentType)
	values.Set(transactionID, transID)
	values.Set(sum, qp.Sum)

	return baseUrl + "?" + values.Encode()
}

func generateTransactionId() string {
	b := make([]byte, 6)
	_, _ = rand.Read(b) // 6 байт = 12 hex символов
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("%d-%s", timestamp, hex.EncodeToString(b))
}
