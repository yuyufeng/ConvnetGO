package convnetlib

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"os/exec"

	"github.com/songgao/water"
	"github.com/songgao/water/waterutil"
)

func GetCvnIP(userid int) net.IP {
	var ip = 0x0A6E0000
	userid++                       //和老版本适配
	offset := (userid / 254) * 2   //每255个地址中.0和.255无法使用
	ip = ip + int(userid) + offset //补位网络地址和广播地址
	data, _ := IntToBytes32(ip, 4) //换算为byte
	self := net.IP(data)           //转换成ip
	return self
}

func Setip() {

	var mask = net.IPv4Mask(255, 0, 0, 0)
	self := GetCvnIP(client.MyUserid)
	Log("myCvnIP:", self)
	client.MyCvnIP = self.String()
	setupIfce(net.IPNet{IP: self, Mask: mask}, client.g_ifce.Name()) //网卡地址绑定
}

func TapInit() {

	if client.g_ifce == nil {

		ifce, err := water.New(water.Config{
			DeviceType: water.TAP,
		})
		if err != nil {
			log.Fatal(err)
		}

		client.mymac = Getmymac(ifce.Name())
		Log("网卡名称:", ifce.Name())
		client.g_ifce = ifce
		defer teardownIfce(ifce)
		dataCh, errCh := startRead(client.g_ifce) //启动网卡

		for { //塞入chain
			select {
			case buffer := <-dataCh:
				tarmac := waterutil.MACDestination(buffer)
				//查找用户，socket

				if bytes.Equal(tarmac, client.mymac) {
					Log("ok")
				}

				if !waterutil.IsBroadcast(waterutil.MACDestination(buffer)) {
					user := GetUserByMac(tarmac)
					if user != nil {
						user.SendBuff(buffer)
						//记录发送字节
						user.Con_send = user.Con_send + int64(len(buffer))
					}
					continue
				} else {
					if buffer[12] == 8 && buffer[13] == 6 {
						for _, v := range client.g_AllUser.Users {
							if v.Con_Status == CON_CONNOK {
								v.SendBuff(buffer)
								//记录发送字节
								v.Con_send = v.Con_send + int64(len(buffer))
							}
						}
					}
				}

				//fmt.Print(waterutil.MACDestination(buffer), ",")
				//Log("received frame:\n", buffer)
				continue
			case err := <-errCh:
				Log("TAP读取错误，请重启程序:", err)
				break
			}
		}
	}
}

const BUFFERSIZE = 1600

func startRead(ifce *water.Interface) (dataChan <-chan []byte, errChan <-chan error) {
	dataCh := make(chan []byte)
	errCh := make(chan error)
	go func() {
		for {
			//很奇怪，这里重新分配内存比固定一块内存所需要的消耗要小
			buffer := make([]byte, BUFFERSIZE)
			n, err := ifce.Read(buffer)
			if err != nil {
				errCh <- err
				break
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

func String2Mac(str string) net.HardwareAddr {
	data, _ := hex.DecodeString(str)
	return data
}
