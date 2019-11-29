package convnetlib

import (
	"github.com/labstack/echo"
)

func formatinput(str string) string {
	return str
}

func SetApi(e *echo.Echo) {
	e.GET("/", welcome)
	e.GET("/login", login)
	e.GET("/groupList", grouplist)
	e.GET("/userList", allUserlist)
	e.GET("/info", clientinfo)
	e.GET("/logout", logout)
}
