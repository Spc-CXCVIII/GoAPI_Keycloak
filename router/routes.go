package router

import (
	v1 "github.com/Spc-CXCVIII/GoAPI_Keycloak/controller/API/V1"
	"github.com/labstack/echo/v4"
)

func Routes(e *echo.Echo) {
	version := e.Group("/v1")
	version.POST("/auth", v1.Auth)
}
