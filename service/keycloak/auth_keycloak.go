package keycloak

import (
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/Spc-CXCVIII/GoAPI_Keycloak/models"
	"github.com/Spc-CXCVIII/GoAPI_Keycloak/service"
)

func Auth_Keycloak(Login *models.Login) (map[string]interface{}, int, error) {
	// call keycloak API and return token
	formData := url.Values{}
	formData.Set("grant_type", os.Getenv("GRANT_TYPE"))
	formData.Set("client_id", os.Getenv("CLIENT_ID"))
	formData.Set("client_secret", os.Getenv("CLIENT_SECRET"))
	formData.Set("username", Login.Username)
	formData.Set("password", Login.Password)
	body := formData.Encode()

	// Create request
	req, err := http.NewRequest("POST", os.Getenv("KEYCLOAK_OPENIDAPI")+"/protocol/openid-connect/token", strings.NewReader(body))
	if err != nil {
		return nil, 500, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, 500, err
	}

	data, err := service.CheckResponse(res)
	if err != nil {
		return nil, 500, err
	}

	return data, res.StatusCode, nil
}
