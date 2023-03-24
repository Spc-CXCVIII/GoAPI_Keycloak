package user

import (
	"github.com/Spc-CXCVIII/GoAPI_Keycloak/database"
	"github.com/Spc-CXCVIII/GoAPI_Keycloak/models"
)

func UserIDFromTokenDB(email string) (interface{}, int, error) {
	user_id := new(models.UserID)
	query := `SELECT id
						FROM users
						WHERE email = ?`
	err := database.DB.QueryRow(query, email).Scan(&user_id.ID)
	if err != nil {
		return 0, 500, err
	}

	return user_id.ID, 200, nil
}
