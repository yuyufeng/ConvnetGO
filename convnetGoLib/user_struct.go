package convnetlib

import (
	"net"
	"strings"
)

const (
	CON_DISCONNECT = iota
	CON_CONNECTING
	CON_CONNOK
)

type User struct {
	Con_AContext *net.UDPConn
	Con_Status   int  //连接状态
	ISOnline     bool //是否在线

	UserID         int
	UserName       string
	AuthorPassword string //访问密码
	MacAddress     string //MAC地址

	MyPeerPort    int
	Con_RetryTime int //尝试重连的次数
	Con_send      int64
	Con_recv      int64
	Con_lastSend  int64
	Needpass      bool
	Con_addr      *net.UDPAddr
}

//确认连接后的更新
func (user *User) RefInfoByPack(conn *net.UDPConn, mac string) {
	user.MacAddress = mac
	user.Con_AContext = conn
	addr := conn.RemoteAddr()
	strs := strings.Split(addr.String(), ":")
	user.Con_addr = &net.UDPAddr{IP: net.ParseIP(strs[0]), Port: Strtoint(strs[1])}
	user.ISOnline = true
}

//刷新用户信息
func (user *User) RefInfoByCmd(ip, port, mac string) {
	user.MacAddress = mac
	user.Con_addr = &net.UDPAddr{IP: net.ParseIP(ip), Port: int(StrToProtocol(port))}
	user.ISOnline = true
}

func UdpSend(conn *net.UDPConn, str string, remoteIP *net.UDPAddr) {
	conn.WriteToUDP([]byte(str), remoteIP)
}

func UdpSendBuff(conn *net.UDPConn, buff []byte, remoteIP *net.UDPAddr) {
	conn.WriteToUDP(buff, remoteIP)
}

//发送信息
func (user *User) SendCmd(message string) {
	UdpSend(user.Con_AContext, message, user.Con_addr)
}

//发送信息
func (user *User) SendBuff(buff []byte) {
	UdpSendBuff(user.Con_AContext, buff, user.Con_addr)
}

func (user *User) TryConnect(userpass string) {

	if userpass != "" {
		user.AuthorPassword = userpass
	}

	if user.Needpass && user.AuthorPassword == "" { //需要密码
		if user.Needpass {
			Log(user, "需要密码")
		}
		user.Con_Status = CON_DISCONNECT
		user.Con_RetryTime = 0
		return
	}

	if user.Con_RetryTime < 7 {
		user.Con_RetryTime++
		sendCmd(ProtocolToStr(cmdCalltoUser) + "," + Inttostr(client.MyUserid) + "," + Inttostr(user.Con_RetryTime) + "," + user.AuthorPassword + "*")
		user.Con_Status = CON_CONNECTING
	} else { //user.Con_RetryTime > 7
		user.Con_RetryTime = 0
		user.Con_Status = CON_DISCONNECT
		return
	}

}
