package subscription

import (
	"context"
	"errors"
	"fmt"
	"net"

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
	db repository.WgPoolRepository
	wg *wireguard.WgClient
}

// Block implements service.SubscriptionService.
func (s *Service) Block(ctx context.Context, wgPeer *models.WgPeer) error {
	return s.wg.BlockPeer(wgPeer.PublicKey)
}

// Enable implements service.SubscriptionService.
func (s *Service) Enable(ctx context.Context, wgPeer *models.WgPeer) error {
	return s.wg.EnablePeer(wgPeer.PublicKey, wgPeer.ProvidedIP)
}

// CreateAccount implements service.SubscriptionService.
func (s *Service) CreateAccount(ctx context.Context, wgPeer *models.WgPeer) (int64, error) {
	allowedIP, priv, pub, err := s.wg.CreateWgPeer()
	if err != nil {
		return 0, fmt.Errorf("cannot create wg peer: %v", err)
	}

	_, provIP, err := net.ParseCIDR(allowedIP)
	if err != nil {
		return 0, fmt.Errorf("cannot parse provided ip: %v", err)
	}

	path, err := s.wg.WriteUserConfig(priv, *provIP)
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
		ServerIP:   net.IPNet{IP: net.ParseIP(s.wg.AddressWithMask())},
		ProvidedIP: *provIP,
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
			ServerIP:   v.ServerIP.String(),
			ProvidedIP: v.ProvidedIP.String(),
			CreatedAt:  v.CreatedAt,
			EndAt:      v.EndAt,
		})
	}

	return &wgPeers, nil
}

func New(db repository.WgPoolRepository, wg *wireguard.WgClient) *Service {

	return &Service{
		db: db,
		wg: wg,
	}
}
