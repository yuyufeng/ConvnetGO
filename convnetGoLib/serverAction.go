package convnetlib

func StartConsole() {

	//创建本地UDP服务
	g_udpport, g_udpserver = SatrtUDPServer(8080, 10)
	//尝试打开UPNP
	UdpServerUpnpSet(g_udpport)
	//创建本地HTTP-API服务
	StartHttpServer(8081, 10)

}
