package payment

import umoney "github.com/Adz-256/cheapVPN/internal/payment/uMoney"

type Payment interface {
	CreatePayLink(qp umoney.Quickpay) (link string, transID string, err error)
}
