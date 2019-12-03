package convnetlib

import (
	"net"

	"github.com/songgao/water"
	"github.com/songgao/water/waterutil"
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
	mymac         net.HardwareAddr

	MyOuterIP string
	MyInnerIp string //内网IP，用于内网互联
	MyCvnIP   string

	g_ifce      *water.Interface
	g_conn      *net.TCPConn
	g_AllUser   *Group
	g_Groups    map[int]*Group
	g_udpserver *net.UDPConn
	g_authtoken string
}

var client Client

func (this Client) logout() {
	//client.g_ifce.Close() //关闭网卡
	client.IsConnectToserver = false
	client.g_conn.Close() //关闭连接
}

func (this Client) writeEther(data []byte) {
	srcmac := waterutil.MACSource(data)
	user := GetUserByMac(srcmac)
	if user != nil {
		if user.Con_Status != CON_CONNOK { //未经过握手则不认为已经连接成功，不允许任何数据进入
			return
		}
		user.Con_recv = user.Con_recv + int64(len(data))
	}

	client.g_ifce.Write(data)

}
