package convnetlib

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os/exec"
	"strings"

	"github.com/songgao/water"
	"github.com/songgao/water/waterutil"
)

func TapInit() {

	if client.g_ifce == nil {

		ifce, err := water.New(water.Config{
			DeviceType: water.TAP,
		})
		if err != nil {
			log.Fatal(err)
		}

		client.Mymac = GetMymac(ifce.Name())

		client.g_ifce = ifce
		defer func() {
			if err := ifce.Close(); err != nil {
				client.g_ifce = nil
				Log(err)
			}
		}()

	}

	var ip = 0x0A000000
	var mask = net.IPv4Mask(255, 255, 255, 0)

	offset := (client.MyUserid / 254) * 2   //每255个地址中.0和.255无法使用
	ip = ip + int(client.MyUserid) + offset //补位网络地址和广播地址

	data, _ := IntToBytes32(ip, 4) //换算为byte
	self := net.IP(data)           //转换成ip

	dataCh, errCh := startRead(client.g_ifce) //启动网卡

	setupIfce(net.IPNet{IP: self, Mask: mask}, client.g_ifce.Name()) //网卡地址绑定

	for { //塞入chain
		select {
		case buffer := <-dataCh:
			tarmac := waterutil.MACDestination(buffer)
			//查找用户，socket

			if bytes.Equal(tarmac, client.Mymac) {
				Log("ok")
				if !waterutil.IsBroadcast(waterutil.MACDestination(buffer)) {
					user := GetUserByMac(tarmac)
					if user != nil {
						user.SendBuff(buffer)
					}
					continue
				}
			}

			fmt.Print(waterutil.MACDestination(buffer), ",")
			//Log("received frame:\n", buffer)
			continue
		case err := <-errCh:
			Log("read error:", err)
		}
	}
}

const BUFFERSIZE = 1522

func startRead(ifce *water.Interface) (dataChan <-chan []byte, errChan <-chan error) {
	dataCh := make(chan []byte)
	errCh := make(chan error)
	go func() {
		for {
			buffer := make([]byte, BUFFERSIZE)
			n, err := ifce.Read(buffer)
			if err != nil {
				errCh <- err
			} else {
				buffer = buffer[:n:n]
				dataCh <- buffer
			}
		}
	}()
	return dataCh, errCh
}

func startPing(dst net.IP, _ bool) {
	if err := exec.Command("ping", "-n", "4", dst.String()).Start(); err != nil {
		Log(err)
	}
}

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

//整形转换成字节
func IntToBytes32(n int, b byte) ([]byte, error) {
	switch b {
	case 1:
		tmp := int8(n)
		bytesBuffer := bytes.NewBuffer([]byte{})
		binary.Write(bytesBuffer, binary.BigEndian, &tmp)
		return bytesBuffer.Bytes(), nil
	case 2:
		tmp := int16(n)
		bytesBuffer := bytes.NewBuffer([]byte{})
		binary.Write(bytesBuffer, binary.BigEndian, &tmp)
		return bytesBuffer.Bytes(), nil
	case 3, 4:
		tmp := int32(n)
		bytesBuffer := bytes.NewBuffer([]byte{})
		binary.Write(bytesBuffer, binary.BigEndian, &tmp)
		return bytesBuffer.Bytes(), nil
	}
	return nil, fmt.Errorf("IntToBytes b param is invaild")
}
