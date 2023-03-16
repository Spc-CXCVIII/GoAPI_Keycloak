package main

import (
	"fmt"

	"github.com/labstack/echo/v4"

	"github.com/Spc-CXCVIII/GoAPI_Keycloak/database"
	"github.com/Spc-CXCVIII/GoAPI_Keycloak/router"
)

func main() {
	// Echo instance
	e := echo.New()
	router.Routes(e)

	// Connect to database
	db, err := database.Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", "1234")))
}
