package gate

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
)

type Gate struct {
}

func New() *Gate {
	return &Gate{}
}

func (g *Gate) Start() {
	go httpStart()
}

func httpStart() {
	loadConfig()

	// 启动Http服务
	e := echo.New()

	// 启动配置热更新
	e.GET("/reload", reload)

	e.Run(standard.New(":8907"))
}
