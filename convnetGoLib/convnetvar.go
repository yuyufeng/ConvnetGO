package convnetlib

import (
	"net"
)

var (
	g_conn *net.TCPConn

	g_serverip          string
	g_serverport        string
	g_serverUdpPort     int
	g_isconnecttoserver bool
	g_MyNatType         int
	g_HasUpnpUDPServer  bool
	g_AllUser           []User
	g_Groups            []Group
	g_udpserver         *net.UDPConn
	g_udpport           int
	g_Userid            int
	g_MyMac             string
	g_MyIp              string
	g_MyInnerIp         string //内网IP，用于内网互联
)
