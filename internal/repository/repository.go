package repository

import (
	"context"

	repoModels "github.com/Adz-256/cheapVPN/internal/repository/psql/models"
)

type UserRepository interface {
	GetAll(ctx context.Context) ([]repoModels.User, error)
	GetUser(ctx context.Context, id int64) (*repoModels.User, error)
	CreateUser(ctx context.Context, user *repoModels.User) (int64, error)
	DeleteUser(ctx context.Context, id int64) error
}

type WgPoolRepository interface {
	GetExpiredAccounts(ctx context.Context) (*[]repoModels.WgPeer, error)
	GetAccountByPublicKey(ctx context.Context, publicKey string) (*repoModels.WgPeer, error)
	GetUserAccounts(ctx context.Context, userID int64) (*[]repoModels.WgPeer, error)
	CreateAccount(ctx context.Context, wgPeer *repoModels.WgPeer) (int64, error)
	UpdateAccount(ctx context.Context, wgPeer *repoModels.WgPeer) error
	DeleteAccount(ctx context.Context, id int64) error
}

type PlanRepository interface {
	GetOneByID(ctx context.Context, id int64) (*repoModels.Plan, error)
	GetAll(ctx context.Context) (*[]repoModels.Plan, error)
	GetAllByCounty(ctx context.Context, county string) (*[]repoModels.Plan, error)
}

type PaymentRepository interface {
	Create(ctx context.Context, payment *repoModels.Payment) (id int64, err error)
	Update(ctx context.Context, payment *repoModels.Payment) error
	Get(ctx context.Context, transID string) (*repoModels.Payment, error)
}
