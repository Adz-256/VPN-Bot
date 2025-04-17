package wg

import (
	"testing"
)

func TestNew(t *testing.T) {
	New("127.0.0.1", "8080", "test")
}
