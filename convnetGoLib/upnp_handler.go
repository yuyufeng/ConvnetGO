package convnetlib

import (
	"context"
	"time"

	"gitlab.com/NebulousLabs/go-upnp"
)

func UdpServerUpnpSet(udpport int) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	d, err := upnp.DiscoverCtx(ctx)
	if err != nil {
		Log(err)
		return err
	}

	// discover external IP
	ip, err := d.ExternalIP()
	if err != nil {
		Log(err)
		return err
	}
	Log("Your external IP is:", ip)

	// forward a port
	err = d.Forward(9001, "upnp test")
	if err != nil {
		Log(err)
	}

	// check that port 9001 is now forwarded
	forwarded, err := d.IsForwardedUDP(uint16(udpport))
	if err != nil {
		Log(err)
	} else if !forwarded {
		Log("port ", udpport, " was not reported as forwarded")
	}

	// un-forward a port
	err = d.Clear(9001)
	if err != nil {
		Log(err)
	}

	// check that port 9001 is no longer forwarded
	forwarded, err = d.IsForwardedTCP(9001)
	if err != nil {
		Log(err)
	} else if forwarded {
		Log("port ", udpport, " should no longer be forwarded")
	}

	// record router's location
	loc := d.Location()
	if err != nil {
		Log(err)
	}
	Log("Loc:", loc)

	// connect to router directly
	d, err = upnp.Load(loc)
	if err != nil {
		Log(err)
	}
	return nil
}
