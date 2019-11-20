package convnetlib

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

func StartHttpServer(port, maxarea int) {

	httptport := port
	//绑定HTTP-API服务

	// Start server
	for {
		e := echo.New()
		SetApi(e)
		err := e.Start("0.0.0.0:" + ProtocolToStr(httptport))
		if err != nil {
			httptport++
			if httptport > port+maxarea {
				return
			}
		} else {
			break
		}
	}

}

func welcome(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
func userlist(c echo.Context) error {
	data, err := json.Marshal(g_Groups)
	if err != nil {
		panic(err)
	}

	if string(data) != "null" {
		return c.String(http.StatusOK, string(data))
	} else {
		return c.String(http.StatusOK, string("{}"))
	}

}
func login(c echo.Context) error {
	//http://127.0.0.1:1323/login?serverip=sh.convnet.net&serverport=23&pass=firefoxinfo&username=yuyuhaso
	username := formatinput(c.QueryParam("username"))
	pass := formatinput(c.QueryParam("pass"))
	g_serverip = formatinput(c.QueryParam("serverip"))
	g_serverport = formatinput(c.QueryParam("serverport"))

	err := ConnectServer(g_serverip, g_serverport)
	if err != nil {
		return c.String(http.StatusOK, "error"+err.Error())
	} else {
		g_isconnecttoserver = true
	}
	g_serverip = strings.Split(g_conn.RemoteAddr().String(), ":")[0]
	g_MyInnerIp = GetPulicIP(g_serverip + ":" + g_serverport)

	//登录请求
	sendCmd(ProtocolToStr(cmdLogin) + "," + username + "," + pass + ",00FFAC539CB9*")
	return c.String(http.StatusOK, "command sent ok")
}
