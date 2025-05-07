package subscription

import (
	"context"
	"errors"
	"fmt"
	"github.com/Adz-256/cheapVPN/internal/config"
	"log/slog"
	"time"

	"github.com/Adz-256/cheapVPN/internal/models"
	"github.com/Adz-256/cheapVPN/internal/repository"
	repoModels "github.com/Adz-256/cheapVPN/internal/repository/psql/models"
	"github.com/Adz-256/cheapVPN/internal/service"
	"github.com/Adz-256/cheapVPN/internal/wireguard"
)

var _ service.SubscriptionService = (*Service)(nil)

var (
	ErrEmptyUserID = errors.New("empty user id")
)

type Service struct {
	db              repository.WgPoolRepository
	wg              *wireguard.WgClient
	updateRateHours int64
}

// Block implements service.SubscriptionService.
func (s *Service) Block(ctx context.Context, pubKey string) error {
	return s.wg.BlockPeer(pubKey)
}

// Enable implements service.SubscriptionService.
func (s *Service) Enable(ctx context.Context, pubkey string) error {
	wg, err := s.db.GetAccountByPublicKey(ctx, pubkey)
	if err != nil {
		return err
	}

	return s.wg.EnablePeer(pubkey, wg.ProvidedIP)
}

// CreateAccount implements service.SubscriptionService.
func (s *Service) CreateAccount(ctx context.Context, wgPeer *models.WgPeer) (int64, error) {
	allowedIP, priv, pub, err := s.wg.CreateWgPeer()
	if err != nil {
		return 0, fmt.Errorf("cannot create wg peer: %v", err)
	}

	path, err := s.wg.WriteUserConfig(priv, allowedIP)
	if err != nil {
		return 0, fmt.Errorf("cannot write user config: %v", err)
	}

	if wgPeer.UserID == 0 {
		return 0, ErrEmptyUserID
	}

	wgPeerRepo := repoModels.WgPeer{
		UserID:     wgPeer.UserID,
		PublicKey:  pub,
		ConfigFile: path,
		ServerIP:   s.wg.AddressWithMask(),
		ProvidedIP: allowedIP,
		EndAt:      wgPeer.EndAt,
	}
	return s.db.CreateAccount(ctx, &wgPeerRepo)
}

// DeleteAccount implements service.SubscriptionService.
func (s *Service) DeleteAccount(ctx context.Context, id int64) error {
	panic("unimplemented")
}

// GetUserAccounts implements service.SubscriptionService.
func (s *Service) GetUserAccounts(ctx context.Context, userID int64) (*[]models.WgPeer, error) {
	var wgPeers []models.WgPeer

	wgPeer, err := s.db.GetUserAccounts(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("cannot get user accounts: %v", err)
	}

	for _, v := range *wgPeer {
		wgPeers = append(wgPeers, models.WgPeer{
			ID:         v.ID,
			PublicKey:  v.PublicKey,
			ConfigFile: v.ConfigFile,
			ServerIP:   v.ServerIP,
			ProvidedIP: v.ProvidedIP,
			CreatedAt:  v.CreatedAt,
			EndAt:      v.EndAt,
		})
	}

	return &wgPeers, nil
}

func (s *Service) GetExpiredAccounts(ctx context.Context) (*[]models.WgPeer, error) {
	var wgPeers []models.WgPeer

	wgPeersRepo, err := s.db.GetExpiredAccounts(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot get expired accounts: %v", err)
	}

	for _, v := range *wgPeersRepo {
		wgPeers = append(wgPeers, models.WgPeer{
			ID:         v.ID,
			PublicKey:  v.PublicKey,
			ConfigFile: v.ConfigFile,
			ServerIP:   v.ServerIP,
			ProvidedIP: v.ProvidedIP,
			CreatedAt:  v.CreatedAt,
			EndAt:      v.EndAt,
		})
	}

	return &wgPeers, nil
}

func (s *Service) StartExpireCRON() {
	t := time.NewTicker(time.Duration(s.updateRateHours) * time.Second)
	defer t.Stop()
	for {
		<-t.C
		accs, err := s.GetExpiredAccounts(context.Background())
		if err != nil {
			slog.Error("GetExpiredAccounts error", slog.Any("error", err))
		}

		for _, acc := range *accs {
			err = s.Block(context.Background(), acc.PublicKey)
			if err != nil {
				slog.Error("Block error", slog.Any("error", err))
			}
		}
	}
}

func NewService(db repository.WgPoolRepository, wg *wireguard.WgClient, cfg config.SubscriptionConfig) *Service {
	return &Service{
		db:              db,
		wg:              wg,
		updateRateHours: cfg.UpdateRateHours(),
	}
}
