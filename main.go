package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"

	"github.com/labstack/echo/v4/middleware"

	"github.com/Spc-CXCVIII/GoAPI_Keycloak/router"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())

	router.Routes(e)
}
