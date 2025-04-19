package fsm

import (
	"context"
)

type Interface interface {
	SetState(ctx context.Context, userID string, state string) string
	GetState(ctx context.Context, userID string) string
}
