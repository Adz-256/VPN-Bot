package psql

import (
	"context"
	"testing"

	repoModels "github.com/Adz-256/cheapVPN/internal/repository/psql/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

func TestPayments(t *testing.T) {
	ctx := context.Background()
	db, err := pgxpool.Connect(ctx, "postgres://postgres:postgres@localhost:5432/postgres")
	if err != nil {
		t.Fatal(err)
		return
	}

	p := NewPayments(db)

	pay := &repoModels.Payment{
		UserID: 1,
		PlanID: 1,
		Amount: 1000,
		Method: "card",
	}

	id, err := p.Create(ctx, pay)
	if err != nil {
		t.Fatal(err)
		return
	}

	pay.ID = id
	pay.Amount = 2000

	err = p.Update(ctx, pay)
	if err != nil {
		t.Fatal(err)
		return
	}

}
