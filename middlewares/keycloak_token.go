package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

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

		// Call the next handler
		return next(c)
	}
}
