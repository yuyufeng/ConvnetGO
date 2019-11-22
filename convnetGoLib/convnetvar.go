package convnetlib

import (
	"net"
)

type Client struct {
	IsConnectToserver bool
	ServerIP          string
	ServerPort        string
	ServerUdpPort     int

	HasUpnpUDP    bool
	UdpServerPort int
	MyNatType     int
	MyUserid      int64
	Mymac         string
	MyOuterIP     string
	MyInnerIp     string //内网IP，用于内网互联

	g_conn      *net.TCPConn
	g_AllUser   *Group
	g_Groups    map[int64]*Group
	g_udpserver *net.UDPConn
}

var client Client
