package helpers

import (
	"github.com/golang-jwt/jwt"
	"github.com/soyandrestrujillo/advanced_go_rest_websockets/models"
	"github.com/soyandrestrujillo/advanced_go_rest_websockets/server"
	"net/http"
	"strings"
)

// GetJWTAuthorizationInfo returns the token and error
func GetJWTAuthorizationInfo(s server.Server, w http.ResponseWriter, r *http.Request) (*jwt.Token, error) {
	tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
	token, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.Config().JWTSecret), nil
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return nil, err
	}

	return token, nil
}
