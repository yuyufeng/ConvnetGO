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
	MyUserid      int
	Mymac         string
	MyOuterIP     string
	MyInnerIp     string //内网IP，用于内网互联

	g_conn      *net.TCPConn
	g_AllUser   Group
	g_Groups    []Group
	g_udpserver *net.UDPConn
}

var client Client
