package wireguard

import (
	"testing"
)

func TestNew(t *testing.T) {
	wg := New("wg1", "127.0.0.1", "51820", "config", "config")
	err := wg.Init()
	if err != nil {
		t.Fatal(err)
	} //only with root privileges
	wg.lastCreatedIP = "127.0.0.2/24"
	ip, priv, pub, err := wg.CreateWgPeer()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(pub)
	err = wg.BlockPeer(pub)
	if err != nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatal(err)
	}
	_, err = wg.WriteUserConfig(priv, ip)
	if err != nil {
		t.Fatal(err)
	}
	ip, priv, _, err = wg.CreateWgPeer()
	if err != nil {
		t.Fatal(err)
	}

	_, err = wg.WriteUserConfig(priv, ip)
	if err != nil {
		t.Fatal(err)
	}

	err = wg.EnablePeer("3mx73uh02rnp4V+S3Xc14C57JNmNZZXWqP22FdaUowQ=", "127.0.0.2/24")
	if err != nil {
		t.Fatal(err)
	}
}
