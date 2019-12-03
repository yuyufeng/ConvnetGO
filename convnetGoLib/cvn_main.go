package convnetlib

func StartConsole() {
	//CLIENT STRUCT
	//
	//				UPNP PORT			 <———————————————————————————————————				<<<<----Udp direct communication
	//					|													|
	//				UDPSERVER												|					—————————————
	//					|													|				   | Peer client |
	//	client ————>g_Allusers———————————————————							|					—————————————
	//		|								    |							|
	//		|—————— g_Groups——————Group——————users——————>newUDPport<—————port knocking		<<<<----Peer client UDPport
	//		|									|							|
	//		|									|			 ———————————————————————————
	//		|						Other Convnet User————>	| Convnet ShakeHand Server  |	<<<<----TCP server trans communication
	//		|									|			 ———————————————————————————
	//		|									|							|
	//	Http Api(echo server)  ——>Login call peer opation————>JSON response
	//															|
	//														  UI layer

	client.MyUserid = 0
	client.g_AllUser = NewGroup()
	client.g_Groups = make(map[int]*Group)
	client.g_AllUser.GroupName = "用户列表"

	//创建本地UDP服务
	client.UdpServerPort, client.g_udpserver = SatrtUDPServer(8080, 10)
	//尝试打开UPNP
	go UdpServerUpnpSet(client.UdpServerPort)
	//检查UPNP是否可用

	go TapInit()
	//创建本地HTTP-API服务
	StartHttpServer(8082, 10)
}
