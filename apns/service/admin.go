package service

import "github.com/labstack/echo"

func reload(c echo.Context) error {
	loadConfig(false)

	return nil
}
