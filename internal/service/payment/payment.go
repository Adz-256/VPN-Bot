package payment

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Adz-256/cheapVPN/internal/broker"
	"github.com/Adz-256/cheapVPN/internal/webhook"
	"github.com/Adz-256/cheapVPN/internal/webhook/smee"

	"github.com/Adz-256/cheapVPN/internal/models"
	paymentImpl "github.com/Adz-256/cheapVPN/internal/payment"
	umoney "github.com/Adz-256/cheapVPN/internal/payment/uMoney"
	"github.com/Adz-256/cheapVPN/internal/repository"
	repoModels "github.com/Adz-256/cheapVPN/internal/repository/psql/models"
	"github.com/Adz-256/cheapVPN/internal/service"
)

var _ service.PaymentService = (*Service)(nil)

type Service struct {
	db        repository.PaymentRepository
	paymentWH webhook.Webhook
	payment   paymentImpl.Payment
	consumer  broker.Consumer
}

const (
	paid     = "paid"
	canceled = "canceled"
)

var (
	ErrPaymentCanceled        = errors.New("payment canceled")
	ErrPaymentAlreadyApproved = errors.New("payment already approved")
	ErrPaymentAlreadyCanceled = errors.New("payment already canceled")
)

func NewService(db repository.PaymentRepository, paymentWH webhook.Webhook, paymentImpl paymentImpl.Payment, consumer broker.Consumer) *Service {
	return &Service{db: db, paymentWH: paymentWH, payment: paymentImpl, consumer: consumer}
}

// ApprovePayment implements service.PaymentService.
func (s *Service) ApprovePayment(ctx context.Context, transID string) (err error) {
	p, err := s.db.Get(ctx, transID)
	if err != nil {
		return fmt.Errorf("cannot get payment: %v", err) // TODO:
	}

	if p.Status == canceled {
		return ErrPaymentCanceled
	} else if p.Status == paid {
		return ErrPaymentAlreadyApproved
	}

	p.Status = paid
	p.PaidAt.Time = time.Now()
	err = s.db.Update(ctx, p)
	if err != nil {
		return fmt.Errorf("cannot update payment: %v", err)
	}

	return nil
}

// CancelPayment implements service.PaymentService.
func (s *Service) CancelPayment(ctx context.Context, transID string) error {
	p, err := s.db.Get(ctx, transID)
	if err != nil {
		return fmt.Errorf("cannot get payment: %v", err) // TODO:
	}

	if p.Status == canceled {
		return ErrPaymentAlreadyCanceled
	} else if p.Status == paid {
		return ErrPaymentAlreadyApproved
	}
	p.Status = canceled
	err = s.db.Update(ctx, p)
	if err != nil {
		return fmt.Errorf("cannot update payment: %v", err)
	}

	return nil
}

func (s *Service) Get(ctx context.Context, transID string) (*models.Payment, error) {
	p, err := s.db.Get(ctx, transID)
	if err != nil {
		return nil, fmt.Errorf("cannot get payment: %v", err)
	}

	return &models.Payment{
		ID:        p.ID,
		TransID:   p.TransID,
		UserID:    p.UserID,
		PlanID:    p.PlanID,
		Amount:    p.Amount,
		Method:    p.Method,
		Status:    p.Status,
		CreatedAt: p.CreatedAt,
		PaidAt:    p.PaidAt,
	}, nil
}

// CreatePayLink implements service.PaymentService.
func (s *Service) CreatePayLink(ctx context.Context, user models.User, plan models.Plan, receiver string) (link string, transID string, err error) {

	ump := umoney.Quickpay{
		Receiver:     receiver,
		QuickpayForm: "shop",
		Targets:      plan.Country,
		PaymentType:  "SB",
		Sum:          plan.Price.String(),
	}

	link, transID, err = s.payment.CreatePayLink(ump)
	if err != nil {
		return "", "", fmt.Errorf("cannot create payment: %v", err)
	}

	payment := &repoModels.Payment{
		UserID:  user.UserID,
		TransID: transID,
		PlanID:  plan.ID,
		Amount:  plan.Price,
		Method:  "uMoney",
	}

	slog.Info("create payment", slog.Any("payment", payment))

	_, err = s.db.Create(ctx, payment)
	if err != nil {
		return "", "", fmt.Errorf("cannot create payment: %v", err)
	}

	return link, transID, nil
}

func (s *Service) StartPaymentsApprover() {
	slog.Debug("starting payments approver")
	for {
		msg, err := s.consumer.Read(context.TODO())

		if msg == nil {
			slog.Debug("message is nil")
			continue
		}
		if err != nil {
			slog.Error("cannot read payment approver msg", slog.Any("err", err))
			continue
		}
		slog.Debug("get New Message From Broker", slog.Any("msg", msg))
		var notif smee.Notification
		err = json.Unmarshal(msg.Value, &notif)
		slog.Debug("Get New Notification", notif)
		if err != nil {
			slog.Error("cannot map to notification", slog.Any("error", err))
			continue
		}
		if notif.Unaccepted != "true" {
			err = s.ApprovePayment(context.Background(), notif.Label)
			if err != nil {
				slog.Error("cannot approve p", slog.Any("error", err))
				continue
			}
		}
	}
}
