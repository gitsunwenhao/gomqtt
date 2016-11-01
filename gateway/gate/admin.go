package gate

import "github.com/labstack/echo"

func reload(c echo.Context) error {
	loadConfig()
	return nil
}
