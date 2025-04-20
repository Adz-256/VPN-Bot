package psql

import (
	"context"
	"testing"
	"time"

	repoModels "github.com/Adz-256/cheapVPN/internal/repository/psql/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

func TestWgPools(t *testing.T) {
	db, err := pgxpool.Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/postgres")
	if err != nil {
		t.Fatal(err)
	}

	pool := NewWgPools(db)

	id, err := pool.CreateAccount(context.Background(), &repoModels.WgPeer{
		UserID:     1234124124,
		PublicKey:  "123121",
		ConfigFile: "test",
		ServerIP:   "127.0.0.1",
		ProvidedIP: "127.0.0.1",
		EndAt:      time.Now(),
	})
	if err != nil {
		t.Fatal(err)
	}

	acc, err := pool.GetUserAccounts(context.Background(), id)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(acc)

	err = pool.UpdateAccount(context.Background(), &repoModels.WgPeer{
		ID: id,
	})
	if err != nil {
		t.Fatal(err)
	}

	err = pool.DeleteAccount(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
}
