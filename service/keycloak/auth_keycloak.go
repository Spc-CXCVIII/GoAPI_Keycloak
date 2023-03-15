package keycloak

import (
	"errors"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/Spc-CXCVIII/GoAPI_Keycloak/models"
	"github.com/Spc-CXCVIII/GoAPI_Keycloak/service"
)

func Auth_Keycloak(Login *models.Login) (map[string]interface{}, int, error) {
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

	res_data, status, err = service.DoRequest(req)
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
	res_data, status, err = service.DoRequest(req)
	if err != nil {
		return res_data, status, err
	}

	return res_data, status, nil
}
