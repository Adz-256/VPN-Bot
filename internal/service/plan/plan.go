package plan

import (
	"context"

	"github.com/Adz-256/cheapVPN/internal/models"
	"github.com/Adz-256/cheapVPN/internal/repository"
	"github.com/Adz-256/cheapVPN/internal/service"
)

var _ service.PlanService = (*Service)(nil)

type Service struct {
	db repository.PlanRepository
}

func NewService(db repository.PlanRepository) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) GetAll(ctx context.Context) (*[]models.Plan, error) {
	plans, err := s.db.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var plansModel []models.Plan

	for _, plan := range *plans {
		plansModel = append(plansModel, models.Plan{
			ID:           plan.ID,
			Country:      plan.Country,
			DurationDays: plan.DurationDays,
			Price:        plan.Price,
			Description:  plan.Description,
		})
	}

	return &plansModel, nil
}

func (s *Service) GetOneByID(ctx context.Context, id int64) (*models.Plan, error) {
	plan, err := s.db.GetOneByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &models.Plan{
		ID:           plan.ID,
		Country:      plan.Country,
		DurationDays: plan.DurationDays,
		Price:        plan.Price,
		Description:  plan.Description,
	}, nil
}

func (s *Service) GetAllByCounty(ctx context.Context, cntry string) (*[]models.Plan, error) {
	plans, err := s.db.GetAllByCounty(ctx, cntry)
	if err != nil {
		return nil, err
	}

	var plansModel []models.Plan

	for _, plan := range *plans {
		plansModel = append(plansModel, models.Plan{
			ID:           plan.ID,
			Country:      plan.Country,
			DurationDays: plan.DurationDays,
			Price:        plan.Price,
			Description:  plan.Description,
		})
	}

	return &plansModel, nil
}
