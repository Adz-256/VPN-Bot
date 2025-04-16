package psql

import (
	"context"
	"testing"

	repoModels "github.com/Adz-256/cheapVPN/internal/repository/psql/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

func TestPayment(t *testing.T) {
	ctx := context.Background()
	db, err := pgxpool.Connect(ctx, "postgres://postgres:postgres@localhost:5432/postgres")
	if err != nil {
		t.Fatal(err)
	}

	u := NewUsers(db)

	user := &repoModels.User{
		ChatID:   3,
		Username: "test",
	}

	id, err := u.CreateUser(ctx, user)
	if err != nil {
		t.Fatal(err)
	}

	newUser, err := u.GetUser(ctx, id)
	if err != nil {
		t.Fatal(err)
	}

	if newUser.ID != id || newUser.Username != user.Username || newUser.IsAdmin != false || newUser.CreatedAt == user.CreatedAt {
		t.Fatalf("expected %v, got %v", user, newUser)
	}

	err = u.DeleteUser(ctx, id)
	if err != nil {
		t.Fatal(err)
	}

	users, err := u.GetAll(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(users) != 0 {
		t.Fatalf("expected 0 users, got %d", len(users))
	}

}
