package gate

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
)

func adminStart() {
	e := echo.New()

	// configuration hot update
	e.GET("/reload", reload)

	e.Run(standard.New(":8907"))
}

func reload(c echo.Context) error {
	loadConfig(false)

	return nil
}
