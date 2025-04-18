package smee

import "testing"

func Test_wh_Run(t *testing.T) {
	wh := New("", "3000")
	wh.Run()
}
