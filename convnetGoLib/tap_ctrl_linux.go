package convnetlib

import (
	"net"
	"os/exec"

	"github.com/songgao/water"
)

func setupIfce(ipNet net.IPNet, dev string) {

	if err := exec.Command("ip", "link", "set", dev, "up").Run(); err != nil {
		Log(err)
	}
	if err := exec.Command("ip", "addr", "add", ipNet.String(), "dev", dev).Run(); err != nil {
		Log(err)
	}
}

func teardownIfce(ifce *water.Interface) {
	client.g_ifce = nil
	if err := ifce.Close(); err != nil {
		Log(err)
	}
	if err := exec.Command("ip", "link", "set", ifce.Name(), "down").Run(); err != nil {
		Log(err)
	}
}
