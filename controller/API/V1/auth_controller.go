package v1

import (
	"github.com/Spc-CXCVIII/GoAPI_Keycloak/database"
	"github.com/Spc-CXCVIII/GoAPI_Keycloak/models"
	"github.com/Spc-CXCVIII/GoAPI_Keycloak/service/keycloak"
	"github.com/labstack/echo/v4"
)

func Auth(c echo.Context) error {
	// Bind data from request
	user := new(models.Login)
	if err := c.Bind(user); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid request payload"})
	}

	// Call keycloak service API
	res_data, status, err := keycloak.Auth_Keycloak(user)
	if err != nil {
		return c.JSON(status, map[string]string{"error": err.Error()})
	}

	// Check if API return data but error
	if _, ok := res_data["error"]; ok {
		return c.JSON(status, res_data)
	}

	// Get userID and user's name from database
	user_data := new(models.UserDataInToken)
	query := `SELECT id, CONCAT(users.f_name, ' ', users.l_name) AS name
						FROM users
						WHERE email = ?`
	err = database.DB.QueryRow(query, user.Username).Scan(&user_data.ID, &user_data.Name)
	if err != nil {
		return c.JSON(500, map[string]string{"error": err.Error()})
	}

	// Return token and user's data
	return c.JSON(status, map[string]interface{}{
		"token":   res_data,
		"user_id": user_data.ID,
		"name":    user_data.Name,
	})
}
