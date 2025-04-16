package psql

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
)

func TestPlans(t *testing.T) {
	ctx := context.Background()
	db, err := pgxpool.Connect(ctx, "postgres://postgres:postgres@localhost:5432/postgres")
	if err != nil {
		t.Fatal(err)
		return
	}

	p := NewPlans(db)

	plans, err := p.GetAll(ctx)
	if err != nil {
		t.Fatal(err)
		return
	}

	id := (*plans)[0].ID

	t.Log(plans)

	plan, err := p.GetOneByID(ctx, id)

	if err != nil {
		t.Fatal(err)
	}

	t.Log(plan)
}
