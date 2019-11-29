package convnetlib

import (
	"net"

	"github.com/songgao/water"
)

type Client struct {
	IsConnectToserver bool
	ServerIP          string
	ServerPort        string
	ServerUdpPort     int

	HasUpnpUDP    bool
	UdpServerPort int
	MyNatType     int
	MyUserid      int
	Mymac         net.HardwareAddr

	MyOuterIP string
	MyInnerIp string //内网IP，用于内网互联

	g_ifce      *water.Interface
	g_conn      *net.TCPConn
	g_AllUser   *Group
	g_Groups    map[int]*Group
	g_udpserver *net.UDPConn
}

var client Client

func (this Client) logout() {
	//client.g_ifce.Close() //关闭网卡
	client.IsConnectToserver = false
	client.g_conn.Close() //关闭连接
}
