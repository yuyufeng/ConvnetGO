package convnetlib

func StartConsole() {
	client.MyUserid = 0
	client.g_AllUser = NewGroup()
	client.g_Groups = make(map[int]*Group)
	client.g_AllUser.GroupName = "用户列表"

	//创建本地UDP服务
	client.UdpServerPort, client.g_udpserver = SatrtUDPServer(8080, 10)
	//尝试打开UPNP
	go UdpServerUpnpSet(client.UdpServerPort)
	go TapInit()
	//创建本地HTTP-API服务
	StartHttpServer(8082, 10)
}
