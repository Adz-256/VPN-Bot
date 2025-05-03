package psql

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Adz-256/cheapVPN/internal/repository"
	repoModels "github.com/Adz-256/cheapVPN/internal/repository/psql/models"
	"github.com/Adz-256/cheapVPN/utils"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
)

var _ repository.PaymentRepository = (*Payments)(nil)

type Payments struct {
	db *pgxpool.Pool
	b  sq.StatementBuilderType
}

const (
	tablePayments   = "payments"
	idColumn        = "id"
	userIDColumn    = "user_id"
	planIDColumn    = "plan_id"
	amountColumn    = "amount"
	statusColumn    = "status"
	methodColumns   = "method"
	createdAtColumn = "created_at"
	paidAtColumn    = "paid_at"
	transIDColumn   = "transaction_id"
)

var (
	ErrEmptyPaymentID = errors.New("payment id is empty")
)

func (p *Payments) Get(ctx context.Context, transID string) (*repoModels.Payment, error) {
	if transID == "" {
		return nil, ErrEmptyPaymentID
	}

	query, args, err := p.b.Select(
		idColumn,
		transIDColumn,
		userIDColumn,
		planIDColumn,
		amountColumn,
		statusColumn,
		methodColumns,
		createdAtColumn,
		paidAtColumn,
	).From(tablePayments).Where(sq.Eq{transIDColumn: transID}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("cannot build sql query: %v", err)
	}

	var payment repoModels.Payment
	err = p.db.QueryRow(ctx, query, args...).Scan(
		&payment.ID,
		&payment.TransID,
		&payment.UserID,
		&payment.PlanID,
		&payment.Amount,
		&payment.Status,
		&payment.Method,
		&payment.CreatedAt,
		&payment.PaidAt,
	)
	if err != nil {
		return nil, fmt.Errorf("cannot execute sql query: %v", err)
	}

	return &payment, nil
}

func (p *Payments) Create(ctx context.Context, payment *repoModels.Payment) (id int64, err error) {
	mPay, err := utils.StructToMap(payment, true)
	if err != nil {
		return 0, fmt.Errorf("malformed struct: %w", err)
	}

	query, args, err := p.b.Insert(tablePayments).SetMap(mPay).Suffix("RETURNING id").ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build query: %w", err)
	}
	slog.Info("query", slog.String("query", query), slog.Any("args", args))
	err = p.db.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %w", err)
	}

	return id, nil
}

func (p *Payments) Update(ctx context.Context, pay *repoModels.Payment) error {
	mPay, err := utils.StructToMap(pay, false)
	if err != nil {
		return fmt.Errorf("malformed struct: %v", err)
	}

	if pay.ID == 0 {
		return ErrEmptyPaymentID
	}

	query, args, err := p.b.Update(tablePayments).Where(sq.Eq{idColumn: pay.ID}).SetMap(mPay).ToSql()
	if err != nil {
		return fmt.Errorf("cannot build sql query: %v", err)
	}

	_, err = p.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("cannot execute sql query: %v", err)
	}

	return nil
}

func NewPayments(db *pgxpool.Pool) *Payments {
	sq := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	return &Payments{
		db: db,
		b:  sq,
	}
}
