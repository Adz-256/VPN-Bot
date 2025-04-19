package wireguard

import (
	"net"
	"testing"
)

func TestNew(t *testing.T) {
	wg := New("wg0", "127.0.0.1", "51820", "config/wg0.conf", "config")
	// err := wg.Init()
	// if err != nil {
	// 	t.Fatal(err)
	// } //only with root privileges
	wg.lastCreatedIP = "127.0.0.2/24"
	_, priv, pub, err := wg.CreateWgPeer()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(pub)
	err = wg.BlockPeer(pub)
	if err != nil {
		t.Fatal(err)
	}
	ip := net.IPNet{IP: net.ParseIP("127.0.0.2"), Mask: net.CIDRMask(24, 32)}
	path, err := wg.WriteUserConfig(priv, ip)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(path)
	err = wg.EnablePeer("3mx73uh02rnp4V+S3Xc14C57JNmNZZXWqP22FdaUowQ=", "127.0.0.2/24")
	if err != nil {
		t.Fatal(err)
	}
}
