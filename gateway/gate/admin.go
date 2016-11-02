package gate

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
)

func adminStart() {
	// 启动Http服务
	e := echo.New()

	// 启动配置热更新
	e.GET("/reload", reload)

	e.Run(standard.New(":8907"))
}

func reload(c echo.Context) error {
	loadConfig(false)

	return nil
}
