package env

import (
	"os"
	"strconv"
)

type Subscription struct {
	updateRateHours int64
}

func NewSubscription() *Subscription {
	updateRateHours, err := strconv.ParseInt(os.Getenv("UPDATE_RATE"), 10, 64)
	if err != nil {
		updateRateHours = 24
	}
	return &Subscription{updateRateHours}
}

func (s *Subscription) UpdateRateHours() int64 {
	return s.updateRateHours
}
