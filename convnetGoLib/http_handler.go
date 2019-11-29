package convnetlib

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

func StartHttpServer(port, maxarea int) {
	httptport := port

	// Start server
	for {

		e := echo.New()
		e.HideBanner = true

		//绑定HTTP-API服务
		SetApi(e)
		err := e.Start("0.0.0.0:" + ProtocolToStr(httptport))
		if err != nil {
			Log("端口被占用，重启服务")
			httptport++
			if httptport > port+maxarea {
				return
			}
		} else {
			Log("APIU服务已启动于：", httptport)
			break
		}
	}

}

func welcome(c echo.Context) error {
	return c.String(http.StatusOK, "{\"info\":\"Convnet Api\"}")
}

func clientinfo(c echo.Context) error {
	return c.String(http.StatusOK, ToJson(client))
}

func allUserlist(c echo.Context) error {
	return c.String(http.StatusOK, ToJson(client.g_AllUser))
}

func ToJson(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	if string(data) != "null" {
		return string(data)
	} else {
		return "{}"
	}
}

func grouplist(c echo.Context) error {
	return c.String(http.StatusOK, ToJson(client.g_Groups))
}
func logout(c echo.Context) error {
	client.logout()
	return c.String(http.StatusOK, "command sent ok")
}

func login(c echo.Context) error {
	//http://127.0.0.1:1323/login?serverip=sh.convnet.net&serverport=23&pass=asdasd&username=yuyuhaso
	username := formatinput(c.QueryParam("username"))
	pass := formatinput(c.QueryParam("pass"))
	client.ServerIP = formatinput(c.QueryParam("serverip"))
	client.ServerPort = formatinput(c.QueryParam("serverport"))

	err := ConnectServer(client.ServerIP, client.ServerPort)
	if err != nil {
		return c.String(http.StatusOK, "error"+err.Error())
	} else {
		client.IsConnectToserver = true
	}

	client.ServerIP = strings.Split(client.g_conn.RemoteAddr().String(), ":")[0]
	client.MyInnerIp = GetPulicIP(client.ServerIP + ":" + client.ServerPort)

	//登录请求
	sendCmd(ProtocolToStr(cmdLogin) + "," + username + "," + pass + "," + Mymacstr() + "*")
	return c.String(http.StatusOK, "command sent ok")
}
