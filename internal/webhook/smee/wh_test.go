package smee

import (
	"fmt"
	"testing"
)

func Test_wh_Run(t *testing.T) {
	wh := New("", "3000")
	ch := make(chan map[string]any, 1024)
	wh.Run(ch)
	for v := range ch {
		fmt.Println(v)
	}
}
