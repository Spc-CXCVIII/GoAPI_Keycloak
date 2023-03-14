package router

import (
	v1 "github.com/Spc-CXCVIII/GoAPI_Keycloak/controller/API/V1"
	"github.com/labstack/echo/v4"
)

func Routes(e *echo.Echo) {
	version := e.Group("/api/v1")
	version.POST("/login", v1.Auth)
}
