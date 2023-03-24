package keycloak

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
	"strings"

	cus_cache "github.com/Spc-CXCVIII/GoAPI_Keycloak/cache"
	"github.com/Spc-CXCVIII/GoAPI_Keycloak/models"
	"github.com/Spc-CXCVIII/GoAPI_Keycloak/service/requests"
)

func LoginKC(Login *models.Login) (map[string]interface{}, int, error) {
	// check email is verified
	res_data, status, err := ManageUserToken()
	if err != nil {
		return nil, status, err
	}

	// get user's id
	req, err := http.NewRequest("GET", os.Getenv("KEYCLOAK_ADMINAPI")+"/users?email="+Login.Username, nil)
	if err != nil {
		return nil, 500, err
	}
	req.Header.Set("Authorization", "Bearer "+res_data["access_token"].(string))

	res_data, status, err = requests.DoRequest(req)
	if err != nil {
		return res_data, status, err
	}

	if res_data["emailVerified"] == false {
		error_data := make(map[string]interface{})
		error_data["error"] = "Your email not verify"

		return error_data, 401, errors.New("Error")
	}

	// call keycloak API and return token
	formData := url.Values{}
	formData.Set("grant_type", os.Getenv("GRANT_TYPE"))
	formData.Set("client_id", os.Getenv("CLIENT_ID"))
	formData.Set("client_secret", os.Getenv("CLIENT_SECRET"))
	formData.Set("username", Login.Username)
	formData.Set("password", Login.Password)
	body := formData.Encode()

	// Create request
	req, err = http.NewRequest("POST", os.Getenv("KEYCLOAK_OPENIDAPI")+"/protocol/openid-connect/token", strings.NewReader(body))
	if err != nil {
		return nil, 500, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send request
	res_data, status, err = requests.DoRequest(req)
	if err != nil {
		return res_data, status, err
	}

	return res_data, status, nil
}

func LogoutKC(id_token string) (map[string]interface{}, int, error) {
	// Get session id from cache
	data, ok := cus_cache.GetData(id_token)
	if !ok {
		error_data := make(map[string]interface{})
		error_data["error"] = "Can't get session id from cache"
		return error_data, 500, errors.New("Error")
	}
	session_id := data.(map[string]interface{})["sid"].(string)

	res_data, status, err := ManageUserToken()
	if err != nil {
		return nil, status, err
	}

	// new request
	req, err := http.NewRequest("DELETE", os.Getenv("KEYCLOAK_ADMINAPI")+"/sessions/"+session_id, nil)
	if err != nil {
		return nil, 500, err
	}
	req.Header.Set("Authorization", "Bearer "+res_data["access_token"].(string))

	// send request
	res_data, status, err = requests.DoRequest(req)
	if err != nil {
		return res_data, status, err
	}

	return res_data, status, nil
}

func RegisterKC(Register *models.Register) (map[string]interface{}, int, error) {
	res_data, status, err := ManageUserToken()
	if err != nil {
		return nil, status, err
	}

	register_data := map[string]interface{}{
		"email":         Register.Email,
		"emailVerified": false,
		"firstName":     Register.Firstname,
		"lastName":      Register.Lastname,
		"username":      Register.Email,
		"enabled":       true,
		"credentials": []map[string]interface{}{
			{
				"type":      "password",
				"value":     Register.Password,
				"temporary": false,
			},
		},
	}
	register_data_byte, err := json.Marshal(register_data)
	if err != nil {
		return nil, 500, err
	}

	// new request
	req, err := http.NewRequest("POST", os.Getenv("KEYCLOAK_ADMINAPI")+"/users", bytes.NewBuffer(register_data_byte))
	if err != nil {
		return nil, 500, err
	}
	req.Header.Set("Authorization", "Bearer "+res_data["access_token"].(string))
	req.Header.Set("Content-Type", "application/json")

	// send request
	res_data, status, err = requests.DoRequest(req)
	if err != nil {
		return res_data, status, err
	}

	return nil, 201, nil
}
