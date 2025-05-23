package service

import (
	"context"

	"github.com/Adz-256/cheapVPN/internal/models"
)

type PaymentService interface {
	Get(ctx context.Context, transID string) (*models.Payment, error)
	CreatePayLink(ctx context.Context, user models.User, plan models.Plan, reciver string) (link string, transID string, err error)
	ApprovePayment(ctx context.Context, transID string) error
	CancelPayment(ctx context.Context, transID string) error
	StartPaymentsApprover()
}

type UserService interface {
	Create(ctx context.Context, user *models.User) (int64, error)
	GetUser(ctx context.Context, id int64) (*models.User, error)
}

type SubscriptionService interface {
	GetExpiredAccounts(ctx context.Context) (*[]models.WgPeer, error)
	GetUserAccounts(ctx context.Context, userID int64) (*[]models.WgPeer, error)
	CreateAccount(ctx context.Context, wgPeer *models.WgPeer) (int64, error)
	DeleteAccount(ctx context.Context, id int64) error
	Block(ctx context.Context, pubKey string) error
	Enable(ctx context.Context, pubKey string) error
	StartExpireCRON()
}

type PlanService interface {
	GetAll(ctx context.Context) (*[]models.Plan, error)
	GetOneByID(ctx context.Context, id int64) (*models.Plan, error)
	GetAllByCounty(ctx context.Context, cntry string) (*[]models.Plan, error)
}
