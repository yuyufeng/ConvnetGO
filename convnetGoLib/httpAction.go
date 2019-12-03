package convnetlib

import (
	"github.com/labstack/echo"
)

func formatinput(str string) string {
	return str
}

func SetApi(e *echo.Echo) {
	e.GET("/", welcome)             //welcome
	e.GET("/login", login)          //用户登录 serverip serverport pass username
	e.GET("/logout", logout)        //登出
	e.GET("/groupList", grouplist)  //获取组用户
	e.GET("/userList", allUserlist) //获取所有用户
	e.GET("/info", clientinfo)      //本地服务状态
	e.GET("/connectUser", connectUser)
}
