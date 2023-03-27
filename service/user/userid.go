package user

import (
	"github.com/Spc-CXCVIII/GoAPI_Keycloak/database"
)

func UserIDFromTokenDB(email string) (interface{}, int, error) {
	var user_id int
	query := `SELECT id
						FROM users
						WHERE email = ?`
	err := database.DB.QueryRow(query, email).Scan(&user_id)
	if err != nil {
		return 0, 500, err
	}

	return user_id, 200, nil
}
