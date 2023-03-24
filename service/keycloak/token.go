package keycloak

import (
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/Spc-CXCVIII/GoAPI_Keycloak/service/requests"
)

func ManageUserToken() (map[string]interface{}, int, error) {
	// call keycloak API and return token
	formData := url.Values{}
	formData.Set("grant_type", os.Getenv("GRANT_TYPE"))
	formData.Set("client_id", os.Getenv("CLIENT_ID"))
	formData.Set("client_secret", os.Getenv("CLIENT_SECRET"))
	formData.Set("username", os.Getenv("MANAGE_USER_USERNAME"))
	formData.Set("password", os.Getenv("MANAGE_USER_PASSWORD"))
	body := formData.Encode()

	// Create request
	req, err := http.NewRequest("POST", os.Getenv("KEYCLOAK_OPENIDAPI")+"/protocol/openid-connect/token", strings.NewReader(body))
	if err != nil {
		return nil, 500, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send request
	res_data, status, err := requests.DoRequest(req)
	if err != nil {
		return res_data, status, err
	}

	return res_data, status, nil
}
