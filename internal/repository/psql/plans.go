package psql

import (
	"context"
	"fmt"

	"github.com/Adz-256/cheapVPN/internal/repository"
	repoModels "github.com/Adz-256/cheapVPN/internal/repository/psql/models"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
)

var _ repository.PlanRepository = (*Plans)(nil)

type Plans struct {
	db *pgxpool.Pool
	b  sq.StatementBuilderType
}

const (
	plansTable        = "plans"
	countryColumn     = "country"
	durationColumn    = "duration_days"
	priceColumn       = "price"
	descriptionColumn = "description"
)

func (p *Plans) GetAll(ctx context.Context) (*[]repoModels.Plan, error) {

	var plans []repoModels.Plan

	query, ars, err := p.b.Select(idColumn,
		countryColumn,
		durationColumn,
		priceColumn,
		descriptionColumn).From(plansTable).ToSql()

	if err != nil {
		return nil, fmt.Errorf("cannot configure sql query: %v", err)
	}

	rows, err := p.db.Query(ctx, query, ars...)
	if err != nil {
		return nil, fmt.Errorf("cannot execute sql query: %v", err)
	}

	for rows.Next() {
		var plan repoModels.Plan

		err = rows.Scan(&plan.ID, &plan.Country, &plan.DurationDays, &plan.Price, &plan.Description)
		if err != nil {
			return nil, fmt.Errorf("cannot scan row: %v", err)
		}
		plans = append(plans, plan)
	}
	return &plans, nil
}

func (p *Plans) GetOneByID(ctx context.Context, id int64) (*repoModels.Plan, error) {

	var plan repoModels.Plan

	query, ars, err := p.b.Select(idColumn,
		countryColumn,
		durationColumn,
		priceColumn,
		descriptionColumn).From(plansTable).Where(sq.Eq{idColumn: id}).ToSql()

	if err != nil {
		return nil, fmt.Errorf("cannot configure sql query: %v", err)
	}

	err = p.db.QueryRow(ctx, query, ars...).Scan(&plan.ID, &plan.Country, &plan.DurationDays, &plan.Price, &plan.Description)

	if err != nil {
		return nil, fmt.Errorf("cannot execute sql query: %v", err)
	}

	return &plan, nil
}

func (p *Plans) GetAllByCounty(ctx context.Context, county string) (*[]repoModels.Plan, error) {

	var plans []repoModels.Plan

	query, ars, err := p.b.Select(idColumn,
		countryColumn,
		durationColumn,
		priceColumn,
		descriptionColumn).From(plansTable).Where(sq.Eq{countryColumn: county}).OrderBy(durationColumn).ToSql()

	if err != nil {
		return nil, fmt.Errorf("cannot configure sql query: %v", err)
	}

	rows, err := p.db.Query(ctx, query, ars...)
	if err != nil {
		return nil, fmt.Errorf("cannot execute sql query: %v", err)
	}

	for rows.Next() {
		var plan repoModels.Plan

		err = rows.Scan(&plan.ID, &plan.Country, &plan.DurationDays, &plan.Price, &plan.Description)
		if err != nil {
			return nil, fmt.Errorf("cannot scan row: %v", err)
		}
		plans = append(plans, plan)
	}

	return &plans, nil
}

func NewPlans(db *pgxpool.Pool) *Plans {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &Plans{
		db: db,
		b:  b,
	}
}
