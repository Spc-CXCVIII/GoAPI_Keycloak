package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	cus_cache "github.com/Spc-CXCVIII/GoAPI_Keycloak/cache"
	"github.com/Spc-CXCVIII/GoAPI_Keycloak/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// Define a custom middleware function that validates JWT tokens
func ValidateToken(next echo.HandlerFunc) echo.HandlerFunc {
	// Define your RS256 public key as a string
	var publicKey = fmt.Sprintf("-----BEGIN PUBLIC KEY-----\n%s\n-----END PUBLIC KEY-----", os.Getenv("KEYCLOAK_REALM_PUBLIC_KEY"))

	return func(c echo.Context) error {
		authorization := strings.Split(c.Request().Header.Get("Authorization"), " ")
		var tokenString string
		error_message := make(map[string]interface{})
		if len(authorization) == 0 {
			error_message["error"] = "requried token"
			return c.JSON(http.StatusUnauthorized, error_message)
		} else {
			tokenString = authorization[1]
		}

		if tokenString == "" {
			error_message["error"] = "Unauthorized"
			return c.JSON(http.StatusUnauthorized, error_message)
		}

		// Parse the JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Verify that the signing method is RS256
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, errors.New("unexpected signing method: " + token.Header["alg"].(string))
			}

			// Parse the public key
			publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
			if err != nil {
				return nil, errors.New("error parsing public key: " + err.Error())
			}

			// Return the public key to use for verification
			return publicKey, nil
		})
		if err != nil {
			error_message["error"] = err.Error()
			return c.JSON(http.StatusUnauthorized, error_message)
		}

		// Validate the token
		if !token.Valid {
			error_message["error"] = "Unauthorized"
			return c.JSON(http.StatusUnauthorized, error_message)
		}

		claims, ok := token.Claims.(jwt.MapClaims) // data from token
		if !ok {
			error_message["error"] = "Unable to extract claims from token"
			return c.JSON(http.StatusUnauthorized, error_message)
		}

		claimsMap := make(map[string]interface{})
		for key, value := range claims {
			claimsMap[key] = value
		}

		id_req := new(models.UserIDRequest)
		if err := c.Bind(id_req); err != nil {
			return c.JSON(400, map[string]string{"error": "Invalid request payload"})
		}

		// set user's id from request. Can use only in this request
		c.Set("id_req", id_req.ID)
		c.Set("email", claims["email"])
		c.Set("id_token", claims["sub"])

		data, status, err := CheckUserID(claims["email"].(string), id_req.ID)
		if err != nil {
			if err.Error() == "Error" {
				return c.JSON(status, map[string]string{"error": data["error"].(string)})
			} else {
				return c.JSON(status, map[string]string{"error": err.Error()})
			}
		}

		// cache data
		cus_cache.AddData(claims["sub"].(string), claimsMap)

		// Call the next handler
		return next(c)
	}
}
