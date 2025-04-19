package fsm

import (
	"github.com/Adz-256/cheapVPN/internal/fsm"
	"github.com/Adz-256/cheapVPN/internal/service"
)

var _ service.FSMService = (*Service)(nil)

type Service struct {
	fsm.Interface
}

func New(fsm fsm.Interface) *Service {
	return &Service{
		Interface: fsm,
	}
}
