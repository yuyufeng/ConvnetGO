package convnetlib

import (
	"fmt"
	"net"
	"os/exec"
	"strings"

	"github.com/songgao/water"
)

func setupIfce(ipNet net.IPNet, dev string) {
	sargs := fmt.Sprintf("interface ip set address name=REPLACE_ME source=static addr=REPLACE_ME mask=REPLACE_ME gateway=none")
	args := strings.Split(sargs, " ")
	args[4] = fmt.Sprintf("name=%s", dev)
	args[6] = fmt.Sprintf("addr=%s", ipNet.IP)
	args[7] = fmt.Sprintf("mask=%d.%d.%d.%d", ipNet.Mask[0], ipNet.Mask[1], ipNet.Mask[2], ipNet.Mask[3])
	cmd := exec.Command("netsh", args...)
	if err := cmd.Run(); err != nil {
		Log(err)
	}
}

func teardownIfce(ifce *water.Interface) {
	if err := ifce.Close(); err != nil {
		Log(err)
	}
}
