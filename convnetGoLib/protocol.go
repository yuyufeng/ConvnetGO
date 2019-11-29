package convnetlib

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	ALL_DATA = iota //0       //发送数据
	UDP_S2S
	UDP_S2SResp
	UDP_C2S
	UDP_C2SResp
	UDP_C2C
	UDP_C2CResp
	UDP_GETPORT
	UDP_P2PResp
	TCP_C2S
	TCP_C2SResp
	TCP_SvrTrans
	ALL_NOTARRIVE //所有方法无法到达
	NOTCONNECT    //对方无法连接
	DISCONNECT    //断开连接
	SAMEIP_CALL   //相同IP连接
)
const (
	NAT_UNKNOW = iota //0
	NAT_CONE
	NAT_SYMMETRIC
)
const (
	CMDSERVERTRANS = iota //0
	//登录登出
	cmdLogin //1
	cmdLoginResp
	cmdLogout
	cmdLogoutResp
	//更新用户状态
	cmdRenewUserStatus //5
	//更新用户信息
	cmdRenewMYinfo
	cmdRenewMYinforesp
	//获取版本信息
	cmdGetVersionResp
	//获取服务器信息
	cmdGetServerPort
	cmdGetServerPortResp
	//注册
	cmdRegistUser
	cmdRegistUserResp
	//用户消息
	cmdSendMsgtoID //13
	cmdSendMsgtoIDResp
	//获取用户、组信息
	cmdGetFriendInfo //15
	cmdGetFriendInfoResp
	cmdGetGroupInfo
	cmdGetGroupInfoResp
	cmdGetGroupDesc
	cmdGetGroupDescresp
	//获取单独用户信息
	cmdGetUserinfo //21
	cmdGetUserinfoResp
	//上线通知
	cmdOnlinetell //23
	cmdOnlinetellResp
	//下线通知
	cmdOffLinetellResp
	//加入组
	cmdJoinGroup
	cmdJoinGroupResp
	//修改组
	cmdmodifyGroup
	cmdmodifyGroupresp

	//消息无法到达
	cmdMsgcantarrive
	//创建组
	cmdCreateGroup //30
	cmdCreateGroupResp
	//退出组
	cmdQuitGroup
	cmdQuitGroupResp
	//踢出用户
	cmdKickOut
	cmdKickOutResp
	//添加用户
	cmdAddFriend //36
	cmdAddFriendResp

	//删除用户
	cmdDelFriend
	cmdDelFriendResp

	//要求服务器转发数据
	cmdOrdServerTrans
	cmdOrdServerTransResp

	//同意添加
	cmdPeerComfimFriend
	cmdPeerComfimFriendResp

	//拒绝添加
	cmdPeerRefusedFriend //44
	cmdPeerRefusedFriendResp

	//同意加入组
	cmdPeerComfimJoinGroup
	cmdPeerComfimJoinGroupresp

	//要求添加
	cmdPeerOrdFriend
	cmdPeerOrdFriendResp
	//查找用户、组
	cmdFindUser
	cmdFindUserResp
	cmdFindGroup
	cmdFindGroupResp

	//请求连接
	cmdCalltoUser     //55
	cmdCalltoUserResp //56
	cmdCalltoUserNewPort
	cmdCalltoUserNewPortResp
	//端开连接
	cmdDissConnUser
	cmdDissConnUserResp
	//客户端确认端口
	cmdISClientUDP
	cmdP2P          //P2P
	cmdKeeponLine   //心跳包
	cmdUserNeedPass //用户需要密码
	cmdSameipInfo   //相同IP应答
	cmdSameipInforesp

	cmdCheckNatType
	cmdServerSendToClient

	cmdGroupChat
	cmdGroupChatResp
	cmdGroupChatCnatArrive
)

const (
	ConstSaveDataLength = 4
)

func reader(readerChannel chan []byte) {

	for {
		select {
		case data := <-readerChannel:
			cmdField := strings.Split(string(data), ",")
			ExecComand(cmdField)
		}
	}
}

func Log(v ...interface{}) {
	log.Println(v...)
}

func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

//解包
func Unpack(buffer []byte, readerChannel chan []byte) []byte {
	length := len(buffer)

	var i int
	for i = 0; i < length; i = i + 1 {
		if length < i+ConstSaveDataLength {
			break
		}

		messageLength := BytesToInt(buffer[i : i+ConstSaveDataLength])
		if length < i+ConstSaveDataLength+messageLength {
			break
		}
		data := buffer[i+ConstSaveDataLength : i+ConstSaveDataLength+messageLength]
		readerChannel <- data

		i += ConstSaveDataLength + messageLength - 1

	}

	if i == length {
		return make([]byte, 0)
	}
	return buffer[i:]
}

func handleConnection(conn net.Conn) {
	//声明一个临时缓冲区，用来存储被截断的数据
	tmpBuffer := make([]byte, 0)

	//声明一个管道用于接收解包的数据
	readerChannel := make(chan []byte, 16)
	go reader(readerChannel)
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			Log(conn.RemoteAddr().String(), " connection error: ", err)
			return
		}
		tmpBuffer = Unpack(append(tmpBuffer, buffer[:n]...), readerChannel)
	}
}
func Strtoint(intstr string) int {
	i, _ := strconv.ParseInt(intstr, 10, 0)
	return int(i)
}

func Strtoint64(intstr string) int64 {
	i, _ := strconv.ParseInt(intstr, 10, 0)
	return i
}

func Inttostr(intnum int) string {
	return strconv.Itoa(intnum)
}

func ProtocolToStr(protostr int) string {
	return strconv.Itoa(protostr)
}

func StrToProtocol(str string) int {
	i, _ := strconv.ParseInt(str, 10, 0)
	return int(i)
}

func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return int(x)
}

func IntToBytes(i int32) []byte {
	byteBuffer := bytes.NewBuffer([]byte{})
	binary.Write(byteBuffer, binary.BigEndian, i)
	return byteBuffer.Bytes()
}
func sendCmd(str string) {
	sendCmdBuff([]byte(str))
}
func sendCmdBuff(context []byte) {
	client.g_conn.Write(append(append([]byte(""), IntToBytes(int32(len(context)))...), context...))
	// var buffer bytes.Buffer
	// asd := IntToBytes(int32(len(context)))
	// buffer.Write(asd)
	// buffer.Write(context)
	// c.Write(buffer.Bytes())
}
