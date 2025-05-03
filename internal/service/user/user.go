package user

import (
	"context"

	"github.com/Adz-256/cheapVPN/internal/models"
	"github.com/Adz-256/cheapVPN/internal/repository"
	repoModels "github.com/Adz-256/cheapVPN/internal/repository/psql/models"
	"github.com/Adz-256/cheapVPN/internal/service"
)

var _ service.UserService = (*Service)(nil)

type Service struct {
	db repository.UserRepository
}

func NewService(db repository.UserRepository) *Service {
	return &Service{
		db: db,
	}
}

func (u *Service) Create(ctx context.Context, user *models.User) (int64, error) {
	repoUser := repoModels.User{
		ID:        user.ID,
		ChatID:    user.UserID,
		Username:  user.Username,
		IsAdmin:   user.IsAdmin,
		CreatedAt: user.CreatedAt,
	}
	return u.db.CreateUser(ctx, &repoUser)
}

func (u *Service) GetUser(ctx context.Context, id int64) (*models.User, error) {
	repoUser, err := u.db.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}
	return &models.User{
		ID:        repoUser.ID,
		UserID:    repoUser.ChatID,
		Username:  repoUser.Username,
		IsAdmin:   repoUser.IsAdmin,
		CreatedAt: repoUser.CreatedAt,
	}, nil
}
