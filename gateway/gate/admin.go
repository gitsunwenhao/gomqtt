package gate

import (
	"github.com/labstack/echo"
)

func adminStart() {
	e := echo.New()

	// configuration hot update
	e.GET("/reload", reload)

	err := e.Start(":8907")
	if err != nil {
		e.Logger.Fatal(err.Error())
	}
}

func reload(c echo.Context) error {
	loadConfig(false)

	return nil
}
