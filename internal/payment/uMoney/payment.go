package umoney

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/Adz-256/cheapVPN/internal/config"
	"log/slog"
	"net/http"
	"net/url"
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
	umoneyAccount string
}

type Quickpay struct {
	Receiver     string
	QuickpayForm string
	Targets      string
	Label        string
	PaymentType  string
	Sum          string
}

func New(cfg config.PaymentConfig) *Payment {
	return &Payment{
		umoneyAccount: cfg.AccountID(),
	}
}

// https://yoomoney.ru/quickpay/confirm.xml?receiver=4100117034899495&quickpay-form=shop&targets=Sponsor%20this%20project&paymentType=SB&sum=5
func (p *Payment) CreatePayLink(qp Quickpay) (link string, transID string, err error) {
	transID = qp.Label

	if transID == "" {
		transID = generateTransactionId()
	}
	if qp.Receiver == "" {
		qp.Receiver = p.umoneyAccount
	}

	payURL := configureRequestURL(qp, transID)

	slog.Info("CreatePayLink", slog.String("payURL", payURL))
	// Создаем POST-запрос по ссылке
	req, err := http.NewRequest(http.MethodPost, payURL, nil)
	if err != nil {
		// логируем, если нужно
		return "", "", err
	}

	client := &http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error {
		link = req.URL.String()
		return nil
	}}

	resp, err := client.Do(req)
	if err != nil {
		// логируем, если нужно
		return "", "", err
	}
	defer resp.Body.Close()

	return link, transID, nil
}

func configureRequestURL(qp Quickpay, transID string) string {
	values := url.Values{}
	values.Set(receiver, qp.Receiver)
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
