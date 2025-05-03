package umoney

import (
	"testing"
)

func Test_generateSHA(t *testing.T) {
	id1 := generateTransactionId()
	id2 := generateTransactionId()

	if id1 == id2 {
		t.Errorf("generateTransactionId() = %v, want %v", id1, id2)
	}

	t.Logf("id1 = %v, id2 = %v", id1, id2)
}

func Test_configureRequestURL(t *testing.T) {
	type args struct {
		qp      Quickpay
		transID string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "configureRequestURL", args: args{qp: Quickpay{
			Receiver:     "4100117034899495",
			QuickpayForm: "shop",
			Targets:      "Sponsor this project",
			PaymentType:  "SB",
			Sum:          "5",
		}, transID: "123"}, want: "https://yoomoney.ru/quickpay/confirm.xml?receiver=4100117034899495&quickpay-form=shop&targets=Sponsor%20this%20project&paymentType=SB&sum=5"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := configureRequestURL(tt.args.qp, tt.args.transID); got != tt.want {
				t.Errorf("configureRequestURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
