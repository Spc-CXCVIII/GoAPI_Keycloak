package middlewares

import (
	"errors"

	"github.com/Spc-CXCVIII/GoAPI_Keycloak/service/user"
)

func CheckUserID(email string, id_req int) (map[string]interface{}, int, error) {
	id_db, status, err := user.UserIDFromTokenDB(email)
	error_data := make(map[string]interface{})
	if err != nil {
		if err.Error() == "Error" {
			error_data["error"] = "Can't get user's id from database"
			return error_data, status, errors.New("Error")
		} else {
			error_data["error"] = err.Error()
			return error_data, status, err
		}
	}

	if id_db.(int) != id_req {
		error_data["error"] = "User's id doesn't match"
		return error_data, 401, errors.New("Error")
	}

	return map[string]interface{}{"status": true}, status, nil
}
