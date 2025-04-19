package payment

import (
	"context"
	"errors"
	"fmt"

	"github.com/Adz-256/cheapVPN/internal/models"
	paymentImpl "github.com/Adz-256/cheapVPN/internal/payment"
	umoney "github.com/Adz-256/cheapVPN/internal/payment/uMoney"
	"github.com/Adz-256/cheapVPN/internal/repository"
	repoModels "github.com/Adz-256/cheapVPN/internal/repository/psql/models"
	"github.com/Adz-256/cheapVPN/internal/service"
)

var _ service.PaymentService = (*Service)(nil)

type Service struct {
	db repository.PaymentsRepository
	paymentImpl.Payment
}

const (
	approved = "approved"
	canceled = "canceled"
)

var (
	ErrPaymentCanceled        = errors.New("payment canceled")
	ErrPaymentAlreadyApproved = errors.New("payment already approved")
	ErrPaymentAlreadyCanceled = errors.New("payment already canceled")
)

// ApprovePayment implements service.PaymentService.
func (s *Service) ApprovePayment(ctx context.Context, transID string) (err error) {
	p, err := s.db.Get(ctx, transID)
	if err != nil {
		return fmt.Errorf("cannot get payment: %v", err) // TODO:
	}

	if p.Status == canceled {
		return ErrPaymentCanceled
	} else if p.Status == approved {
		return ErrPaymentAlreadyApproved
	}

	p.Status = approved
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
	} else if p.Status == approved {
		return ErrPaymentAlreadyApproved
	}
	p.Status = canceled
	err = s.db.Update(ctx, p)
	if err != nil {
		return fmt.Errorf("cannot update payment: %v", err)
	}

	return nil
}

// CreatePayLink implements service.PaymentService.
func (s *Service) CreatePayLink(ctx context.Context, user models.User, plan models.Plan, reciver string) (link string, err error) {

	ump := umoney.Quickpay{
		Recieiver:    reciver,
		QuickpayForm: "shop",
		Targets:      plan.Name,
		PaymentType:  "SB",
		Sum:          plan.Price.String(),
	}

	link, transID, err := s.Payment.CreatePayLink(ump)
	if err != nil {
		return "", fmt.Errorf("cannot create payment: %v", err)
	}
	_, err = s.db.Create(ctx, &repoModels.Payment{
		UserID:  user.ID,
		TransID: transID,
		PlanID:  plan.ID,
		Status:  "created",
		Amount:  plan.Price,
		Method:  "uMoney",
	})
	if err != nil {
		return "", fmt.Errorf("cannot create payment: %v", err)
	}

	return link, nil
}
