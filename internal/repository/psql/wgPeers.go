package psql

import (
	"context"
	"fmt"
	"time"

	"log/slog"

	"github.com/Adz-256/cheapVPN/internal/repository"
	repoModels "github.com/Adz-256/cheapVPN/internal/repository/psql/models"
	"github.com/Adz-256/cheapVPN/utils"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
)

var _ repository.WgPoolRepository = (*WgPeers)(nil)

type WgPeers struct {
	db *pgxpool.Pool
	b  sq.StatementBuilderType
}

const (
	WgPoolsTable     = "wg_peers"
	publicKeyColumn  = "public_key"
	configFileColumn = "config_file"
	serverIPColumn   = "server_ip"
	providedIPColumn = "provided_ip"
	endAtColumn      = "end_at"
	blockedColumn    = "blocked"
)

// Get all accounts that belongs to user
func (w *WgPeers) GetUserAccounts(ctx context.Context, userID int64) (*[]repoModels.WgPeer, error) {
	query, args, err := w.b.Select(idColumn,
		userIDColumn, publicKeyColumn,
		configFileColumn, serverIPColumn,
		providedIPColumn, createdAtColumn, endAtColumn).From(WgPoolsTable).Where(sq.Eq{userIDColumn: userID}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("cannot build sql query: %v", err)
	}
	rows, err := w.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("cannot execute sql query: %v", err)
	}
	defer rows.Close()

	var wgPeers []repoModels.WgPeer
	for rows.Next() {
		var wgPeer repoModels.WgPeer

		err := rows.Scan(&wgPeer.ID, &wgPeer.UserID, &wgPeer.PublicKey,
			&wgPeer.ConfigFile, &wgPeer.ServerIP, &wgPeer.ProvidedIP,
			&wgPeer.CreatedAt, &wgPeer.EndAt)
		if err != nil {
			return nil, fmt.Errorf("cannot scan row: %v", err)
		}

		wgPeers = append(wgPeers, wgPeer)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	return &wgPeers, nil
}

func (w *WgPeers) CreateAccount(ctx context.Context, wgPeer *repoModels.WgPeer) (int64, error) {
	peerMap, err := utils.StructToMap(wgPeer, true)
	if err != nil {
		return 0, fmt.Errorf("malformed struct: %v", err)
	}

	query, args, err := w.b.Insert(WgPoolsTable).Columns(publicKeyColumn, configFileColumn, serverIPColumn, providedIPColumn).SetMap(peerMap).Suffix("RETURNING id").ToSql()
	if err != nil {
		return 0, fmt.Errorf("cannot build sql query: %v", err)
	}
	slog.Info("query", slog.String("query", query), slog.Any("args", args))
	var id int64
	err = w.db.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("cannot execute sql query: %v", err)
	}

	return id, nil
}

func (w *WgPeers) UpdateAccount(ctx context.Context, wgPeer *repoModels.WgPeer) error {
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

func (w *WgPeers) DeleteAccount(ctx context.Context, id int64) error {
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

func (w *WgPeers) GetAccountByPublicKey(ctx context.Context, pubKey string) (*repoModels.WgPeer, error) {
	query, args, err := w.b.Select(idColumn,
		userIDColumn, publicKeyColumn,
		configFileColumn, serverIPColumn,
		providedIPColumn, createdAtColumn, endAtColumn).From(WgPoolsTable).Where(sq.Eq{publicKeyColumn: pubKey}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("cannot build sql query: %v", err)
	}

	var wgPeer repoModels.WgPeer

	err = w.db.QueryRow(ctx, query, args...).Scan(
		&wgPeer.ID, &wgPeer.UserID, &wgPeer.PublicKey, &wgPeer.ConfigFile,
		&wgPeer.ServerIP, &wgPeer.ProvidedIP, &wgPeer.CreatedAt, &wgPeer.EndAt)

	if err != nil {
		return nil, fmt.Errorf("cannot execute sql query: %v", err)
	}

	return &wgPeer, nil
}

func (w *WgPeers) GetExpiredAccounts(ctx context.Context) (*[]repoModels.WgPeer, error) {
	query, args, err := w.b.Select(idColumn,
		userIDColumn, publicKeyColumn,
		configFileColumn, serverIPColumn,
		providedIPColumn, createdAtColumn, endAtColumn).From(WgPoolsTable).Where(sq.Lt{endAtColumn: time.Now()}, sq.Eq{blockedColumn: false}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("cannot build sql query: %v", err)
	}

	rows, err := w.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("cannot execute sql query: %v", err)
	}
	defer rows.Close()

	var wgPeers []repoModels.WgPeer
	for rows.Next() {
		var wgPeer repoModels.WgPeer
		err = rows.Scan(
			&wgPeer.ID, &wgPeer.UserID, &wgPeer.PublicKey, &wgPeer.ConfigFile,
			&wgPeer.ServerIP, &wgPeer.ProvidedIP, &wgPeer.CreatedAt, &wgPeer.EndAt)
		if err != nil {
			return nil, fmt.Errorf("cannot scan row: %v", err)
		}
		wgPeers = append(wgPeers, wgPeer)
	}

	return &wgPeers, nil
}

func NewWgPools(db *pgxpool.Pool) *WgPeers {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &WgPeers{
		db: db,
		b:  b,
	}
}
