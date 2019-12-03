package convnetlib

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"net"
)

func SatrtUDPServer(localport, maxarea int) (int, *net.UDPConn) {
	var defaultport = localport
	var defaultportstr = "0.0.0.0:" + ProtocolToStr(defaultport)

	for {
		udpAddr, _ := net.ResolveUDPAddr("udp", defaultportstr)
		net.ResolveUDPAddr("udp", defaultportstr)
		conn, err := net.ListenUDP("udp", udpAddr)

		//defer conn.Close()
		if err != nil { //如果绑定失败
			defaultport++ //换个端口
			defaultportstr = ":" + ProtocolToStr(defaultport)
			if defaultport > localport+maxarea {
				return 0, nil
			}
		} else {
			conn.SetReadBuffer(1600)
			conn.SetWriteBuffer(1600)
			go udpProcess(conn)
			return defaultport, conn
		}
	}
}

func GetPortFromServer(port, localport int, serverip string, keepconn bool) (int, *net.UDPConn) {
	//创建本地的udp访问远程端口
	var defaultport = localport
	var conn *net.UDPConn
	var err error
	var remoteip net.IP

	remoteip = net.ParseIP(serverip)

	lAddr := &net.UDPAddr{Port: defaultport}
	rAddr1 := &net.UDPAddr{IP: remoteip, Port: port}

	for {
		conn, err = net.DialUDP("udp", lAddr, rAddr1)
		conn.SetReadBuffer(1600)
		conn.SetWriteBuffer(1600)

		if !keepconn {
			defer conn.Close()
		}
		if err != nil { //如果绑定失败
			defaultport++ //换个端口
			lAddr = &net.UDPAddr{Port: defaultport}
			if defaultport > localport+100 {
				return 0, nil
			}
		} else {
			break
		}
	}
	conn.Write([]byte("a"))
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	buf := make([]byte, 5)
	len, _, _ := conn.ReadFromUDP(buf)
	var portres int
	if len == 0 {
		return 0, conn
	} else {
		portr := string(buf[:len])
		portres, _ = strconv.Atoi(portr)
	}

	return portres, conn
}

func udpProcess(conn *net.UDPConn) {
	for {
		data := make([]byte, BUFFERSIZE)
		n, remoteAddr, err := conn.ReadFromUDP(data)
		//fmt.Println(n)

		if err != nil {
			fmt.Println("failed to read UDP msg because of ", err.Error())
			continue
		}

		if data[0] == '0' && n > 12 {
			client.writeEther(data[2:])
			continue
		}

		cmdField := strings.Split(string(data), ",")
		//UDP服务端接收
		ExecUdpComand(conn, remoteAddr, cmdField)
	}
}

//UDPPACKE:|UDP_P2PResp|UDP_C2C|userid|mac|ordertoken

func ExecUdpComand(conn *net.UDPConn, remoteAddr *net.UDPAddr, cmdField []string) {
	//UDP服务端接收
	switch StrToProtocol(cmdField[0]) {
	//接收数据之前处理了
	case cmdISClientUDP: //验证是否是本地UDP服务端口
		if cmdField[1] == "ConVnetCheck"+GenToken() {
			client.HasUpnpUDP = true
			client.MyNatType = NAT_CONE
		}
	case DISCONNECT:
		cmdGetFriendInfoRespDecode(cmdField)
	case UDP_P2PResp:
		cmdUDP_P2PResp(conn, remoteAddr, cmdField)
	default:
		Log("尚未实现的ExecUdpComand", cmdField)
	}
}

func cmdUDP_P2PResp(conn *net.UDPConn, remoteAddr *net.UDPAddr, cmdField []string) {
	//UDP服务端接收针对于P2P部分的处理
	//Log("Udp resp", cmdField)
	tmpuser := client.g_AllUser.GetUserByid(Strtoint(cmdField[2]))

	Log(tmpuser.UserNickName, "已经成功连接")
	//这种接入方式基本上只要知道了对方的对接端口就可以完成接入申请
	//解释：如果知道了对方的IP和对接端口（很随机了），那么就认可为允许接入目前并无不可，当然确实是有安全隐患
	//TODO:服务器应该为双方的握手行为加上TOKEN校验
	switch StrToProtocol(cmdField[1]) {
	case UDP_C2C, UDP_C2S, UDP_S2S: //================>	握手应答
		tmpuser.RefInfoByUdpPack(conn, remoteAddr, cmdField[3])
		responseportocol := ProtocolToStr(StrToProtocol(cmdField[1]) + 1) //+1(response)
		//收到就成功握手，回应请求完成握手
		//                        CMD                 TYPE         				ID           			   MAC
		tmpstr := ProtocolToStr(UDP_P2PResp) + "," + responseportocol + "," + Inttostr(int(client.MyUserid)) + "," + mymacstr() + ","
		//UDP握手应答
		UdpSend(conn, tmpstr, remoteAddr)

	case UDP_C2CResp, UDP_C2SResp, UDP_S2SResp: //    <================response无需应答
		//收到就成功握手，回应请求完成握手
		tmpuser.RefInfoByUdpPack(conn, remoteAddr, cmdField[3])

	case ALL_NOTARRIVE:
		tmpuser.Con_conType = 0
		tmpuser.Con_Status = CON_DISCONNECT
	default:
		Log("尚未实现的cmdUDP_P2PResp", cmdField)
	}
}
