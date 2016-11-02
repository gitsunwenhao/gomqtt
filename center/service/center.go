package service

import (
	"fmt"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
)

type Center struct {
}

func New() *Center {
	return &Center{}
}

func (c *Center) Start(isStatic bool) {
	fmt.Println(isStatic)
	loadConfig(isStatic)

	go httpStart()
}

func httpStart() {
	// 启动Http服务
	e := echo.New()

	// 启动配置热更新
	e.GET("/reload", reload)

	e.Run(standard.New(":8908"))
}
