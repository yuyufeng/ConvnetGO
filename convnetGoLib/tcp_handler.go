package convnetlib

import (
	"log"

	"net"
	"strconv"
	"strings"
)

func ConnectServer(server string, port string) error {
	log.Println("connect.", server, ":", port)
	var err error
	if client.IsConnectToserver {
		client.g_conn.Close()
	}

	service := server + ":" + port
	tcpAddr, _ := net.ResolveTCPAddr("tcp", service)
	client.g_conn, err = net.DialTCP("tcp", nil, tcpAddr)

	if err != nil {
		return err
	}

	// read or write on conn
	go HandleConn()
	return nil
}

func HandleConn() {
	defer func() {
		client.g_conn.Close()

		client.logout()
		log.Printf(client.ServerIP + ":" + client.ServerPort + "连接断开")
	}()
	handleConnection(client.g_conn)
}

func Split_string(s string) []string {
	a := strings.Split(s, ",")
	return a
}

//获取外网IP
func GetPulicIP(serveruri string) string {
	conn, _ := net.Dial("tcp", serveruri)
	defer conn.Close()
	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")
	return localAddr[0:idx]
}

func ExecComand(cmdField []string) {
	switch StrToProtocol(cmdField[0]) {
	case cmdLoginResp:
		cmdLoginRespDecode(cmdField)
	case cmdGetFriendInfoResp:
		cmdGetFriendInfoRespDecode(cmdField)
	case cmdGetGroupInfoResp:
		cmdGetGroupInfoRespDecode(cmdField)
	case cmdGetServerPortResp:
		cmdGetServerPortRespDecode(cmdField)
	case cmdOnlinetellResp:
		cmdOnlinetellRespDecode(cmdField)
	case cmdCalltoUserResp:
		cmdCalltoUserRespDecode(cmdField)
	case cmdKickOutResp:
		cmdKickOutRespDecode(cmdField)
	case cmdSameipInforesp:
		cmdSameipInforesppDecode(cmdField)
	default:
		Log("尚未实现", cmdField)
	}
}

func CheckNat(port1, port2 string) {
	port1int, _ := strconv.Atoi(port1)
	port2int, _ := strconv.Atoi(port2)

	client.MyNatType = NAT_UNKNOW
	port1int, _ = GetPortFromServer(port1int, 7700, client.ServerIP, false)
	port2int, _ = GetPortFromServer(port2int, 7700, client.ServerIP, false)

	if port1int == 0 || port2int == 0 {
		Log("udpNatType================== NAT_UNKNOW")
		client.MyNatType = NAT_UNKNOW
		return
	}

	if port1int == port2int {
		client.MyNatType = NAT_CONE //CONE NAT 最具穿透力的类型
		Log("udpNatType================== NAT_CONE")
	} else {
		client.MyNatType = NAT_SYMMETRIC //S NAT 有可能可以穿透
		Log("udpNatType================== NAT_SYMMETRIC")
	}
}

func Udpconfim(port string) string {
	serverAddr := client.ServerIP + ":" + port
	conn, err := net.Dial("udp", serverAddr)
	if err != nil {
		Log("Can't resolve address: ", err)
		return ""
	}
	defer conn.Close()

	conn.Write([]byte("a"))
	buf := make([]byte, BUFFERSIZE)
	conn.Read(buf)
	if err != nil {
		return ""
	}
	return string(buf)
}

func cmdOnlinetellRespDecode(cmdField []string) {
	Log("用户上线", cmdField)
	user := client.g_AllUser.GetUserByid(Strtoint(cmdField[2]))
	if user != nil {
		user.TryConnect("")
	}
}

func cmdKickOutRespDecode(cmdField []string) {
	Log("用户离开组", cmdField)
	group := client.g_Groups[Strtoint(cmdField[2])]
	group.RemoveUserByid(Strtoint(cmdField[2]))
}

func mymacstr() string {
	str := client.mymac.String()

	return strings.ToUpper(strings.Replace(str, ":", "", -1))
}

func Getmymac(etherName string) net.HardwareAddr {

	// 获取本机的MAC地址
	interfaces, err := net.Interfaces()
	if err != nil {
		panic("Error : " + err.Error())
	}
	for _, inter := range interfaces {
		//mac := inter.HardwareAddr //获取本机MAC地址
		if etherName == inter.Name {
			//fmt.Println("MAC = ", mac)
			return inter.HardwareAddr
		}
	}

	return nil
}

func cmdSameipInforesppDecode(cmdField []string) {

	tmpuserid := Strtoint(cmdField[1])
	tmpuser := client.g_AllUser.GetUserByid(tmpuserid)
	Log(tmpuser.UserNickName, "相同IP内网呼叫")
	if tmpuser == nil {
		return
	}
	tmpuser.Dissconnect()
	tmpuser.con_addr = &net.UDPAddr{IP: net.ParseIP(cmdField[5]), Port: Strtoint(cmdField[2])}
	tmpstr := ProtocolToStr(UDP_P2PResp) + "," + ProtocolToStr(UDP_S2S) + "," + Inttostr(client.MyUserid) + "," + mymacstr() + ","
	UdpSend(client.g_udpserver, tmpstr, tmpuser.con_addr)
}

func cmdCalltoUserRespDecode(cmdField []string) {
	Log("用户请求连接", cmdField)
	//cmd+连接协议+用户ID+用户IP+用户端口+用户mac
	tmpuserid := Strtoint(cmdField[2])
	tmpuser := client.g_AllUser.GetUserByid(tmpuserid)
	if tmpuser == nil {
		return
	}

	//tmpuser.RefInfoByCmd(cmdField[3], cmdField[4], cmdField[5])
	switch StrToProtocol(cmdField[1]) {

	case SAMEIP_CALL:
		Log(cmdField)
		Log(tmpuser.UserNickName, "呼入方IP和本机相同")
		Log("通知对方", ProtocolToStr(cmdSameipInfo)+","+cmdField[2]+","+ProtocolToStr(client.UdpServerPort)+","+"0"+","+mymacstr()+","+client.MyInnerIp+"*")
		//通知对方使用本地IP进行尝试通讯，跳转到UDP_S2S或者UDPC2S
		sendCmd(ProtocolToStr(cmdSameipInfo) + "," + cmdField[2] + "," + ProtocolToStr(client.UdpServerPort) + "," + "0" + "," + mymacstr() + "," + client.MyInnerIp + "*")

	case UDP_S2S, UDP_C2S:
		Log("呼入方准备好直连")
		tmpstr := ProtocolToStr(UDP_P2PResp) + "," + ProtocolToStr(UDP_S2SResp) + "," + Inttostr(int(client.MyUserid)) + "," + mymacstr() + ","
		UdpSend(client.g_udpserver, tmpstr, tmpuser.con_addr)
		UdpSend(client.g_udpserver, tmpstr, tmpuser.con_addr)
	case UDP_GETPORT:
		//本地不具备UPNP连接的情况下服务器要求本地准备一个临时端口
		//获取临时端口后绑定给用户
		//通知服务器我已准备好，可以尝试握手
		Log("为", tmpuser.UserNickName, "准备本地对撞端口")
		tmpuser.Dissconnect()
		int, conn := GetPortFromServer(client.ServerUdpPort, 10800, client.ServerIP, true)
		tmpuser.con_AContext = conn
		if client.MyNatType == NAT_SYMMETRIC {
			int++ //如果是非对称端口下次通讯端口号至少+1，先尝试到这个程度，一般可以奏效，不行也就只能再做步长猜测了
			//TODO，可以进一步尝试端口号编号步长猜测
		}
		sendCmd(ProtocolToStr(cmdCalltoUserNewPort) + "," + cmdField[2] + "," + Inttostr(int) + "*")
	case TCP_SvrTrans:
		Log(tmpuser.UserNickName, "服务器转发接入")
		//7次打洞全部失败，服务器允许的情况下会建立TCP转发
		tmpuser := client.g_AllUser.GetUserByid(Strtoint(cmdField[2]))
		tmpuser.Con_Status = CON_CONNOK
		tmpuser.Con_conType = 2
		tmpuser.MacAddress = String2Mac(cmdField[3])
	}
}

func cmdGetServerPortRespDecode(cmdField []string) {
	Log("获取udp服务", cmdField[1], cmdField[2])
	CheckNat(cmdField[1], cmdField[2])
	client.ServerUdpPort = Strtoint(cmdField[1])
	var mynat string
	switch client.MyNatType {
	case NAT_UNKNOW:
		mynat = "UK"
	case NAT_UPNP, NAT_CONE:
		mynat = "CN"
	case NAT_SYMMETRIC:
		mynat = "SN"
	}
	//通知服务器本地NAT类型，是否可以upnp直连
	//						cmd							type				udpport					tcpport&endstar
	sendCmd(ProtocolToStr(cmdRenewUserStatus) + "," + mynat + "," + Inttostr(client.ServerUdpPort) + ",0*")
}

func cmdGetFriendInfoRespDecode(cmdField []string) {
	//返回用户信息
	var tmpGroup *Group
	var strstep = 3
	tmpGroup = client.g_Groups[0]
	if tmpGroup == nil {
		tmpGroup = NewGroup()
		tmpGroup.GroupID = 0
		tmpGroup.GroupName = "好友组"
		client.g_Groups[0] = tmpGroup
	} else {
		tmpGroup.ClearUser()
	}

	for i := 0; i < ((len(cmdField)-1)/strstep)-1; i++ {
		tmpuserid := Strtoint(cmdField[i*strstep+1])
		tmpuser := client.g_AllUser.GetUserByid(tmpuserid)

		if tmpuser == nil {
			user := &User{}
			user.UserID = tmpuserid
			user.UserNickName = cmdField[i*strstep+2]
			user.ISOnline = cmdField[i*strstep+3] == "T"
			client.g_AllUser.Adduser(user)

			tmpGroup.Adduser(user)
		} else {
			tmpGroup.Adduser(tmpuser)
		}
	}
}

func cmdGetGroupInfoRespDecode(cmdField []string) {
	var strstep = 4
	var tmpgourpid int
	var tmpGroup *Group
	Log("组信息创建")

	//Log("好友列表", cmdField)
	for i := 0; i < ((len(cmdField)-1)/strstep)-1; i++ {
		if cmdField[i*strstep+1] == "G" {
			tmpgourpid = Strtoint(cmdField[i*strstep+3])
			tmpGroup = client.g_Groups[tmpgourpid]
			if tmpGroup == nil {
				tmpGroup = NewGroup()
				tmpGroup.GroupID = tmpgourpid
				tmpGroup.GroupName = cmdField[i*strstep+2]
				client.g_Groups[tmpgourpid] = tmpGroup
			} else {
				tmpGroup.ClearUser()
			}
		}

		if cmdField[i*strstep+1] == "U" {
			tmpuserid := Strtoint(cmdField[i*strstep+3])
			tmpuser := client.g_AllUser.GetUserByid(tmpuserid)
			if tmpuser == nil {
				user := &User{}
				user.UserID = tmpuserid
				user.UserNickName = cmdField[i*strstep+2]
				user.ISOnline = cmdField[i*strstep+4] == "T"
				client.g_AllUser.Adduser(user)
				tmpGroup.Adduser(user)
			} else {
				tmpGroup.Adduser(tmpuser)
			}
		}
	}

	sendCmd(ProtocolToStr(cmdOnlinetell) + "*")
	Log("通知上线")
}

func cmdUserNeedPassDecode(cmdField []string) {
	user := client.g_AllUser.GetUserByid(Strtoint(cmdField[1]))
	if user != nil {
		user.Needpass = true
		user.AuthorPassword = ""
	}
}

func cmdLoginRespDecode(cmdField []string) { //实现Getname方法
	switch cmdField[1] {
	case "T":
		Log("登录成功", cmdField)
		Log("用户名：", cmdField[4], "用户IP", cmdField[2], "昵称", cmdField[5])
		client.MyOuterIP = cmdField[2]

		client.MyUserid = Strtoint(cmdField[3])
		go CheckUpnp()

		Setip()
		//获取NAT类型辅助确认端口
		sendCmd(ProtocolToStr(cmdGetServerPort) + "*")
		//获取好友列表
		sendCmd(ProtocolToStr(cmdGetFriendInfo) + "*")
		//获取组信息
		sendCmd(ProtocolToStr(cmdGetGroupInfo) + "*")
	case "D":
		Log("重复登录", cmdField)
		client.IsConnectToserver = false
		client.g_conn.Close()
	case "F":
		Log("登录失败", cmdField)
		client.IsConnectToserver = false
		client.g_conn.Close()
	default:
		Log("无输出")
	}
}
