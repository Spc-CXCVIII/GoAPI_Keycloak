package router

import (
	v1 "github.com/Spc-CXCVIII/GoAPI_Keycloak/controller/API/V1"
	"github.com/Spc-CXCVIII/GoAPI_Keycloak/middlewares"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Routes(e *echo.Echo) {
	version := e.Group("/api/v1")
	version.Use(middleware.Logger())

	// auth group
	auth := version.Group("/auth")
	auth.POST("/login", v1.Login)
	auth.Use(middlewares.ValidateToken)
	auth.POST("/logout", v1.Logout)

	// user group
	user := version.Group("/user")
	user.Use(middlewares.ValidateToken)
}
