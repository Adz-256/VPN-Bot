package psql

import (
	"context"
	"fmt"

	"github.com/Adz-256/cheapVPN/internal/repository"
	repoModels "github.com/Adz-256/cheapVPN/internal/repository/psql/models"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Users struct {
	db *pgxpool.Pool
	b  sq.StatementBuilderType
}

var _ repository.UserRepository = (*Users)(nil)

const (
	tableUsers = "users"

	chatIDColumn   = "chat_id"
	usernameColumn = "username"
	isAdminColumn  = "is_admin"
)

func (u *Users) GetAll(ctx context.Context) ([]repoModels.User, error) {
	query, args, err := u.b.Select(
		idColumn,
		chatIDColumn,
		usernameColumn,
		isAdminColumn,
		createdAtColumn,
	).From(tableUsers).ToSql()
	if err != nil {
		return nil, fmt.Errorf("cannot build sql query: %v", err)
	}

	rows, err := u.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("cannot execute sql query: %v", err)
	}
	defer rows.Close()

	var users []repoModels.User
	for rows.Next() {
		var user repoModels.User
		err = rows.Scan(
			&user.ID,
			&user.ChatID,
			&user.Username,
			&user.IsAdmin,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("cannot scan row: %v", err)
		}
		users = append(users, user)
	}

	return users, nil
}

func (u *Users) GetUser(ctx context.Context, id int64) (*repoModels.User, error) {
	query, args, err := u.b.Select(
		idColumn,
		chatIDColumn,
		usernameColumn,
		isAdminColumn,
		createdAtColumn,
	).From(tableUsers).Where(sq.Eq{idColumn: id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("cannot build sql query: %v", err)
	}

	var user repoModels.User
	err = u.db.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.ChatID,
		&user.Username,
		&user.IsAdmin,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("cannot execute sql query: %v", err)
	}

	return &user, nil
}

func (u *Users) CreateUser(ctx context.Context, user *repoModels.User) (int64, error) {
	query, args, err := u.b.Insert(tableUsers).Columns(
		chatIDColumn,
		usernameColumn,
	).Values(
		user.ChatID,
		user.Username,
	).Suffix("RETURNING id").ToSql()
	if err != nil {
		return 0, fmt.Errorf("cannot build sql query: %v", err)
	}

	var id int64
	err = u.db.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("cannot execute sql query: %v", err)
	}

	return id, nil
}

func (u *Users) DeleteUser(ctx context.Context, id int64) error {
	query, args, err := u.b.Delete(tableUsers).Where(sq.Eq{idColumn: id}).ToSql()
	if err != nil {
		return fmt.Errorf("cannot build sql query: %v", err)
	}

	_, err = u.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("cannot execute sql query: %v", err)
	}

	return nil
}

func NewUsers(db *pgxpool.Pool) *Users {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &Users{
		db: db,
		b:  b,
	}
}
