package v1

import (
	"github.com/Spc-CXCVIII/GoAPI_Keycloak/database"
	"github.com/Spc-CXCVIII/GoAPI_Keycloak/models"
	"github.com/Spc-CXCVIII/GoAPI_Keycloak/service/keycloak"

	"github.com/labstack/echo/v4"
)

func Login(c echo.Context) error {
	// Bind data from request
	user := new(models.Login)
	if err := c.Bind(user); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid request payload"})
	}

	// Call keycloak API
	res_data, status, err := keycloak.LoginKC(user)
	if err != nil {
		if err.Error() == "Error" { // Check if API return error
			return c.JSON(status, map[string]string{"error": res_data["error"].(string)})
		} else {
			return c.JSON(status, map[string]string{"error": err.Error()})
		}
	}

	// Get userID and user's name from database
	user_data := new(models.UserDataInToken)
	query := `SELECT id, CONCAT(users.first_name, ' ', users.last_name) AS name
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

func Logout(c echo.Context) error {
	// Check user id
	id_token, ok := c.Get("id_token").(string)
	if !ok {
		return c.JSON(401, map[string]string{"error": "Cache data not found"})
	}

	res_data, status, err := keycloak.LogoutKC(id_token)
	if err != nil {
		if err.Error() == "Error" {
			return c.JSON(status, map[string]string{"error": res_data["error"].(string)})
		} else {
			return c.JSON(status, map[string]string{"error": err.Error()})
		}
	}

	return c.JSON(200, map[string]interface{}{
		"message": "Logout success",
		"status":  "success",
	})
}

func Register(c echo.Context) error {
	// Bind data from request
	regis_data := new(models.Register)
	if err := c.Bind(regis_data); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid request payload"})
	}

	// Check if email already exist
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users
						WHERE email = ? OR id_card = ? OR phone = ?)`
	err := database.DB.QueryRow(query, regis_data.Email, regis_data.IDNumber, regis_data.Phone).Scan(&exists)
	if err != nil {
		return c.JSON(500, map[string]string{"error": err.Error()})
	}

	if exists {
		return c.JSON(409, map[string]string{"error": "Email or ID number or phone number already exist"})
	} else {
		// Add new user to Keycloak
		res_data, status, err := keycloak.RegisterKC(regis_data)
		if err != nil {
			if err.Error() == "Error" { // Check if API return error
				return c.JSON(status, map[string]string{"error": res_data["error"].(string)})
			} else {
				return c.JSON(status, map[string]string{"error": err.Error()})
			}
		}

		// Add new user to database
		var id_or_ps string
		if regis_data.TH {
			id_or_ps = "id_card"
		} else {
			id_or_ps = "passport"
		}
		query := "INSERT INTO users(title, first_name, last_name, email, `" + id_or_ps + "`, phone)\n" +
			"VALUES (?, ?, ?, ?, ?, ?)"

		// Prepare statement
		stmt, err := database.DB.Prepare(query)
		if err != nil {
			return c.JSON(500, map[string]string{"error": err.Error()})
		}
		defer stmt.Close()

		// Execute statement
		_, err = stmt.Exec(regis_data.Title, regis_data.Firstname, regis_data.Lastname, regis_data.Email, regis_data.IDNumber, regis_data.Phone)
		if err != nil {
			return c.JSON(500, map[string]string{"error": err.Error()})
		}

		return c.JSON(201, map[string]interface{}{
			"status":  "success",
			"message": "Register new user success, please check your email to verify your account. It may appear in your spam folder.",
		})
	}
}
