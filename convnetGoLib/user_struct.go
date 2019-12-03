package convnetlib

import (
	"net"
	"strings"
	"time"
)

const (
	CON_DISCONNECT = iota
	CON_CONNECTING
	CON_CONNOK
)

type User struct {
	noCopy       noCopy
	con_AContext *net.UDPConn
	Con_Status   int  //连接状态
	ISOnline     bool //是否在线

	UserID         int
	UserNickName   string
	AuthorPassword string           //访问密码
	MacAddress     net.HardwareAddr //MAC地址

	myPeerPort    int
	con_RetryTime int //尝试重连的次数
	Con_send      int64
	Con_recv      int64
	con_lastSend  int64
	Needpass      bool
	con_addr      *net.UDPAddr
	Con_conType   int //1 udp 2 tcptrans
}

//确认连接后的更新
func (user *User) RefInfoByUdpPack(conn *net.UDPConn, remoteAddr *net.UDPAddr, mac string) {
	user.MacAddress = String2Mac(mac)
	user.con_AContext = conn
	strs := strings.Split(remoteAddr.String(), ":")
	user.con_addr = &net.UDPAddr{IP: net.ParseIP(strs[0]), Port: Strtoint(strs[1])}

	user.Con_Status = CON_CONNOK
	user.Con_conType = 1
	user.ISOnline = true
}

//刷新用户信息
func (user *User) RefInfoByCmd(ip, port, mac string) {

	user.MacAddress = String2Mac(mac)

	user.con_addr = &net.UDPAddr{IP: net.ParseIP(ip), Port: int(StrToProtocol(port))}
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
	UdpSend(user.con_AContext, message, user.con_addr)
}

//发送信息
func (user *User) SendBuff(buff []byte) {
	switch user.Con_conType {
	case 1: //UDP方式发送
		buff = append([]byte("0,"), buff...)
		UdpSendBuff(user.con_AContext, buff, user.con_addr)
	case 2: //服务器转发
		//			转发		myid						targetuserid			split	appendbuff
		tempstr := "0," + Inttostr(client.MyUserid) + "," + Inttostr(user.UserID) + "*"
		buff = append([]byte(tempstr), buff...)
		sendCmdBuff(buff)
	}

}

func (user *User) TryConnect(userpass string) {
	user.Con_Status = CON_DISCONNECT
	user.con_AContext = nil
	user.con_RetryTime = 3
	user.con_addr = nil

RetryConnect:
	if user.Con_Status == CON_CONNOK {
		return
	}

	if userpass != "" {
		user.AuthorPassword = userpass
	}

	if user.Needpass && user.AuthorPassword == "" { //需要密码
		if user.Needpass {
			Log(user, "需要密码")
		}
		user.Con_Status = CON_DISCONNECT
		user.con_RetryTime = 0
		return
	}

	if user.con_RetryTime < 7 {
		user.con_RetryTime++
		sendCmd(ProtocolToStr(cmdCalltoUser) + "," + Inttostr(int(client.MyUserid)) + "," + Inttostr(user.con_RetryTime) + "," + user.AuthorPassword + "*")
		user.Con_Status = CON_CONNECTING
	} else { //user.Con_RetryTime > 7
		user.con_RetryTime = 0
		user.Con_Status = CON_DISCONNECT
		return
	}

	time.Sleep(time.Second * 3)
	goto RetryConnect

}
