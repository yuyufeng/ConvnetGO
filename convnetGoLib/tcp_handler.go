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
		client.IsConnectToserver = false
		log.Printf(client.ServerIP + ":" + client.ServerPort + "连接断开")
	}()
	handleConnection(client.g_conn)
}

func Split_string(s string) []string {
	a := strings.Split(s, ",")
	return a
}

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
	buf := make([]byte, 1024)
	conn.Read(buf)
	if err != nil {
		return ""
	}
	return string(buf)
}

func cmdOnlinetellRespDecode(cmdField []string) {
	Log("用户上线", cmdField)
	user := getUserByid(client.g_AllUser, Strtoint(cmdField[2]))
	user.TryConnect("")
}

func cmdKickOutRespDecode(cmdField []string) {
	Log("用户离开组", cmdField)

}

func cmdCalltoUserRespDecode(cmdField []string) {
	Log("用户请求连接", cmdField)
	//cmd+连接协议+用户ID+用户IP+用户端口+用户mac
	tmpuserid := Strtoint(cmdField[2])
	tmpuser := getUserByid(client.g_AllUser, tmpuserid)

	//cmd+连接协议+用户ID+用户IP+用户端口+用户mac
	tmpuser.RefInfoByCmd(cmdField[3], cmdField[4], cmdField[5])

	switch StrToProtocol(cmdField[1]) {

	case SAMEIP_CALL:
		Log(tmpuser.UserName, "呼入方IP和本机相同")
		sendCmd(ProtocolToStr(cmdSameipInfo) + "," + cmdField[2] + "," + ProtocolToStr(client.UdpServerPort) + "," + "0" + "," + client.Mymac + "," + client.MyInnerIp + "*")
		//CALL_TO_USER_RESP-UDP_S2S

	case UDP_S2S, UDP_C2S:
		Log("呼入方准备好直连")
		tmpstr := ProtocolToStr(UDP_P2PResp) + "," + ProtocolToStr(UDP_S2SResp) + "," + ProtocolToStr(client.MyUserid) + "," + client.Mymac + ","
		UdpSend(client.g_udpserver, tmpstr, tmpuser.Con_addr)
		UdpSend(client.g_udpserver, tmpstr, tmpuser.Con_addr)
		UdpSend(client.g_udpserver, tmpstr, tmpuser.Con_addr)

	case UDP_GETPORT:
		int, conn := GetPortFromServer(client.ServerUdpPort, 10800, client.ServerIP, true)
		tmpuser.Con_AContext = conn
		if client.MyNatType == NAT_SYMMETRIC {
			int++
		}
		sendCmd(ProtocolToStr(cmdCalltoUserNewPort) + "," + cmdField[2] + "," + Inttostr(int) + "*")
	}
}

func cmdGetServerPortRespDecode(cmdField []string) {
	Log("获取udp服务", cmdField[1], cmdField[2])
	CheckNat(cmdField[1], cmdField[2])
	client.ServerUdpPort = Strtoint(cmdField[1])
}

func cmdGetFriendInfoRespDecode(cmdField []string) {
	//返回用户信息
	var tmpGroup *Group
	var strstep = 3
	var tmpuserid int
	tmpGroup = getGroupByid(0)
	if tmpGroup == nil {
		tmpGroup = &Group{}
		tmpGroup.GroupID = 0
		tmpGroup.GroupName = "好友组"
		client.g_Groups = append(client.g_Groups, *tmpGroup)
	} else {
		tmpGroup.ClearUser()
	}
	tmpGroup = getGroupByid(0)
	for i := 0; i < ((len(cmdField)-1)/strstep)-1; i++ {
		tmpuserid, _ = strconv.Atoi(cmdField[i*strstep+1])
		tmpuser := getUserByid(client.g_AllUser, tmpuserid)
		if tmpuser == nil {
			user := &User{}
			user.UserID = tmpuserid
			user.UserName = cmdField[i*strstep+2]
			user.ISOnline = cmdField[i*strstep+3] == "T"
			client.g_AllUser = append(client.g_AllUser, *user)
			tmpGroup.Adduser(user)
		} else {
			tmpGroup.Adduser(tmpuser)
		}
	}
}

func cmdGetGroupInfoRespDecode(cmdField []string) {
	var strstep = 4
	var tmpuserid, tmpgourpid int
	var tmpGroup *Group
	Log("组信息创建")

	//Log("好友列表", cmdField)
	for i := 0; i < ((len(cmdField)-1)/strstep)-1; i++ {
		if cmdField[i*strstep+1] == "G" {
			tmpgourpid, _ = strconv.Atoi(cmdField[i*strstep+3])
			tmpGroup = getGroupByid(tmpgourpid)
			if tmpGroup == nil {
				tmpGroup = new(Group)
				tmpGroup.GroupID = tmpgourpid
				tmpGroup.GroupName = cmdField[i*strstep+2]
				client.g_Groups = append(client.g_Groups, *tmpGroup)
			} else {
				tmpGroup.ClearUser()
			}
		}

		//go语言这里一定要重新获取一下，出了赋值的作用域?
		//应该是垃圾回收机制的问题，不重新get一下的话会毫无作用
		tmpGroup = getGroupByid(tmpgourpid)
		if cmdField[i*strstep+1] == "U" {
			tmpuserid, _ = strconv.Atoi(cmdField[i*strstep+3])
			tmpuser := getUserByid(client.g_AllUser, tmpuserid)
			if tmpuser == nil {
				user := &User{}
				user.UserID = tmpuserid
				user.UserName = cmdField[i*strstep+2]
				user.ISOnline = cmdField[i*strstep+4] == "T"
				client.g_AllUser = append(client.g_AllUser, *user)
				tmpGroup.Adduser(user)
			} else {
				tmpGroup.Adduser(tmpuser)
			}
		}
	}
}

func getGroupByid(tmpgourpid int) *Group {
	for i := 0; i < len(client.g_Groups); i++ {
		if client.g_Groups[i].GroupID == tmpgourpid {
			return &client.g_Groups[i]
		}
	}
	return nil
}
func getUserByid(userlist []User, tmpuserid int) *User {
	for i := 0; i < len(userlist); i++ {
		if userlist[i].UserID == tmpuserid {
			return &userlist[i]
		}
	}
	return nil
}

func cmdUserNeedPassDecode(cmdField []string) {
	user := getUserByid(client.g_AllUser, Strtoint(cmdField[1]))
	if user != nil {
		user.Needpass = true
		user.AuthorPassword = ""
	}
}

func cmdLoginRespDecode(cmdField []string) { //实现Getname方法
	switch cmdField[1] {
	case "T":
		Log("登录成功", cmdField)
		client.MyOuterIP = cmdField[2]
		client.MyUserid = Strtoint(cmdField[3])
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
