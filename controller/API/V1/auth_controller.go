package v1

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func Auth(c echo.Context) error {
	fmt.Println("auth")
	return c.String(200, "auth")
}
