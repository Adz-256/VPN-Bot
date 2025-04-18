package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

type fsm struct {
	client *redis.Client
}

const (
	stateField = "state"
)

//In Redis the data of user will be stored as hash table with userID as key
//and state, loadV, loadW as values

func New(addr string, port string) *fsm {
	fsm := &fsm{
		client: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", addr, port),
			Password: "",
			DB:       0,
		}),
	}
	if fsm.client.Ping(context.Background()).Err() != nil {
		log.Fatal("cannot connect to redis")
	}
	return fsm
}

// SetState sets the state associated with the given userID in Redis.
// It returns an error if something fails while executing the command in Redis.
func (f *fsm) SetState(ctx context.Context, userID string, state string) error {
	return f.client.HSet(ctx, userID, stateField, state).Err()
}

// GetState retrieves the state associated with the given userID from Redis.
// It returns the state as a string.
func (f *fsm) State(ctx context.Context, userID string) string {
	return f.client.HGet(ctx, userID, stateField).Val()
}
