package convnetlib

import (
	"context"
	"testing"
	"time"

	"gitlab.com/NebulousLabs/go-upnp"
)

func Test_cmdGetServerPortRespDecode(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	d, err := upnp.DiscoverCtx(ctx)
	if err != nil {
		t.Skip(err)
	}

	// discover external IP
	ip, err := d.ExternalIP()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Your external IP is:", ip)

	// forward a port
	err = d.Forward(9001, "upnp test")
	if err != nil {
		t.Fatal(err)
	}

	// check that port 9001 is now forwarded
	forwarded, err := d.IsForwardedTCP(9001)
	if err != nil {
		t.Fatal(err)
	} else if !forwarded {
		t.Fatal("port 9001 was not reported as forwarded")
	}

	// un-forward a port
	err = d.Clear(9001)
	if err != nil {
		t.Fatal(err)
	}

	// check that port 9001 is no longer forwarded
	forwarded, err = d.IsForwardedTCP(9001)
	if err != nil {
		t.Fatal(err)
	} else if forwarded {
		t.Fatal("port 9001 should no longer be forwarded")
	}

	// record router's location
	loc := d.Location()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Loc:", loc)

	// connect to router directly
	d, err = upnp.Load(loc)
	if err != nil {
		t.Fatal(err)
	}

}
