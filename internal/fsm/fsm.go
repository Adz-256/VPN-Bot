package fsm

import (
	"context"
)

type Interface interface {
	SetState(ctx context.Context, userID string, state string) string
	State(ctx context.Context, userID string) string
}
