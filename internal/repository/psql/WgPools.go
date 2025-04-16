package psql

import (
	"context"
	"fmt"

	"github.com/Adz-256/cheapVPN/internal/repository"
	repoModels "github.com/Adz-256/cheapVPN/internal/repository/psql/models"
	"github.com/Adz-256/cheapVPN/utils"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
)

var _ repository.WgPoolRepository = (*WgPools)(nil)

type WgPools struct {
	db *pgxpool.Pool
	b  sq.StatementBuilderType
}

const (
	WgPoolsTable     = "wg_pools"
	publicKeyColumn  = "public_key"
	configFileColumn = "config_file"
	serverIPColumn   = "server_ip"
	providedIPColumn = "provided_ip"
)

func (w *WgPools) GetAccount(ctx context.Context, id int) (*repoModels.WgPeer, error) {
	query, args, err := w.b.Select("*").From(WgPoolsTable).Where(sq.Eq{idColumn: id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("cannot build sql query: %v", err)
	}

	var wgPeer repoModels.WgPeer
	err = w.db.QueryRow(ctx, query, args...).Scan(&wgPeer.ID, &wgPeer.UserID, &wgPeer.PublicKey, &wgPeer.ConfigFile, &wgPeer.ServerIP, &wgPeer.ProvidedIP, &wgPeer.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("cannot execute sql query: %v", err)
	}

	return &wgPeer, nil
}

func (w *WgPools) CreateAccount(ctx context.Context, wgPeer *repoModels.WgPeer) (int, error) {
	query, args, err := w.b.Insert(WgPoolsTable).Columns(publicKeyColumn, configFileColumn, serverIPColumn, providedIPColumn).Values(wgPeer.PublicKey, wgPeer.ConfigFile, wgPeer.ServerIP, wgPeer.ProvidedIP).Suffix("RETURNING id").ToSql()
	if err != nil {
		return 0, fmt.Errorf("cannot build sql query: %v", err)
	}

	var id int
	err = w.db.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("cannot execute sql query: %v", err)
	}

	return id, nil
}

func (w *WgPools) UpdateAccount(ctx context.Context, wgPeer *repoModels.WgPeer) error {
	wgMap, err := utils.StructToMap(wgPeer, false)
	if err != nil {
		return fmt.Errorf("malformed struct: %v", err)
	}
	query, args, err := w.b.Update(WgPoolsTable).SetMap(wgMap).Where(sq.Eq{idColumn: wgPeer.ID}).ToSql()
	if err != nil {
		return fmt.Errorf("cannot build sql query: %v", err)
	}

	_, err = w.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("cannot execute sql query: %v", err)
	}

	return nil
}

func (w *WgPools) DeleteAccount(ctx context.Context, id int) error {
	query, args, err := w.b.Delete(WgPoolsTable).Where(sq.Eq{idColumn: id}).ToSql()
	if err != nil {
		return fmt.Errorf("cannot build sql query: %v", err)
	}

	_, err = w.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("cannot execute sql query: %v", err)
	}

	return nil
}

func NewWgPools(db *pgxpool.Pool) *WgPools {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &WgPools{
		db: db,
		b:  b,
	}
}
